package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/internal/validator"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/homepage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/landingpage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/loginpage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/registerpage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/resetpasswordentercodepage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/resetpasswordpage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/resetpasswordsendcodepage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/usersettingspage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/verifyemailpage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/verifyemailsendcodepage"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	pageData := app.newPageData(r)
	if app.isAuthenticated(r) {
		app.render(r.Context(), w, r, homepage.Page(pageData))
		return
	}
	app.render(r.Context(), w, r, landingpage.Page(pageData))
}

func (app *application) userRegister(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	pageData := registerpage.RegisterPageData{
		PageData: basePageData,
	}
	app.render(r.Context(), w, r, registerpage.Page(pageData))
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	pageData := loginpage.LoginPageData{
		PageData: basePageData,
	}
	app.render(r.Context(), w, r, loginpage.Page(pageData))
}

func (app *application) userRegisterPost(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	pageData := registerpage.RegisterPageData{
		PageData: basePageData,
	}
	pageData.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &pageData)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	pageData.CheckFieldError(pageData.Validate.Var(pageData.Email, "required"), "email", "field is required")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.Email, "email"), "email", "field must be a valid email")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.Email, "max=128"), "email", "cannot contain more than 128 characters")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.Username, "required"), "username", "field is required")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.Username, "min=3"), "username", "must contain at least 3 characters")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.Username, "max=64"), "username", "cannot contain more than 64 characters")
	err = pageData.Validate.Var(pageData.Username, "email")
	if err == nil {
		pageData.AddFieldError("username", "field must not be an email")
	}
	pageData.CheckFieldError(pageData.Validate.Var(pageData.Username, "no_whitespace"), "username", "field cannot contain whitespace characters")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.Password, "required"), "password", "field is required")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.Password, "no_whitespace"), "password", "field cannot contain whitespace characters")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.Password, "min=8"), "password", "must contain at least 8 characters")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.Password, "max=64"), "password", "cannot contain more than 64 characters")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.PasswordConfirmation, "required"), "password-confirmation", "field is required")
	pageData.CheckFieldBool(pageData.Password == pageData.PasswordConfirmation, "password", "passwords do not match")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.Terms, "required,eq=on"), "terms", "You must agree to the Terms & Conditions")

	if !pageData.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, registerpage.Page(pageData))
		return
	}

	_, err = app.userService.CreateUser(pageData.Email, pageData.Username, pageData.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			pageData.AddFieldError("email", "Email address is already in use")
			w.WriteHeader(http.StatusUnprocessableEntity)
			app.render(r.Context(), w, r, registerpage.Page(pageData))
			return
		}
		if errors.Is(err, models.ErrDuplicateUsername) {
			pageData.AddFieldError("username", "Username is already in use")
			w.WriteHeader(http.StatusUnprocessableEntity)
			app.render(r.Context(), w, r, registerpage.Page(pageData))
			return
		}
		app.serverError(w, r, err)
	}

	app.sessionManager.Put(r.Context(), "email-to-verify", pageData.Email)
	http.Redirect(w, r, "/user/verify-email", http.StatusSeeOther)
}

func (app *application) userVerifyEmail(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	email := app.sessionManager.PopString(r.Context(), "email-to-verify")
	pageData := verifyemailpage.VerifyEmailPageData{
		PageData: basePageData,
		Email:    email,
	}
	app.render(r.Context(), w, r, verifyemailpage.Page(pageData))
}

func (app *application) userVerifyEmailPost(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	pageData := verifyemailpage.VerifyEmailPageData{
		PageData: basePageData,
	}
	pageData.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &pageData)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	pageData.CheckFieldError(pageData.Validate.Var(pageData.Code, "required"), "code", "Code field is required")

	if !pageData.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, verifyemailpage.Page(pageData))
		return
	}

	ok, err := app.userService.VerifyEmail(pageData.Code)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !ok {
		pageData.AddFieldError("code", "The verification code you entered is invalid or has expired. Please try again.")
		app.render(r.Context(), w, r, verifyemailpage.Page(pageData))
		return
	}
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userVerifyEmailResendCode(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	pageData := verifyemailsendcodepage.VerifyEmailSendCodePageData{
		PageData: basePageData,
	}
	app.render(r.Context(), w, r, verifyemailsendcodepage.Page(pageData))
}

func (app *application) userVerifyEmailResendCodePost(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	pageData := verifyemailsendcodepage.VerifyEmailSendCodePageData{
		PageData: basePageData,
	}
	pageData.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &pageData)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	pageData.CheckFieldError(pageData.Validate.Var(pageData.Email, "required"), "email", "Email field is required")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.Email, "email"), "email", "Email field must be a valid email")

	if !pageData.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, verifyemailsendcodepage.Page(pageData))
		return
	}

	_, err = app.userService.ResendEmailVerificationCode(pageData.Email)

	if err != nil && !errors.Is(err, models.ErrUserDoesNotExist) && !errors.Is(err, services.ErrEmailAlreadyVerified) {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "email-to-verify", pageData.Email)
	http.Redirect(w, r, "/user/verify-email", http.StatusSeeOther)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	pageData := loginpage.LoginPageData{
		PageData: basePageData,
	}
	pageData.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &pageData)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	pageData.CheckFieldError(pageData.Validate.Var(pageData.Username, "required"), "username", "field is required")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.Password, "required"), "password", "field is required")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.Password, "max=64"), "password", "cannot contain more than 64 characters")

	if !pageData.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, loginpage.Page(pageData))
		return
	}

	user, err := app.userService.AuthenticateUser(pageData.Username, pageData.Password)
	if err != nil && errors.Is(err, services.ErrInvalidCredentials) {
		pageData.AddFieldError("password", "invalid credentials")
		w.WriteHeader(http.StatusUnauthorized)
		app.render(r.Context(), w, r, loginpage.Page(pageData))
		return
	}

	if err != nil && errors.Is(err, services.ErrEmailNotVerified) {
		app.sessionManager.Put(r.Context(), "email-to-verify", user.Email)
		http.Redirect(w, r, "/user/verify-email", http.StatusSeeOther)
		return
	}

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserId", user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserId")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userResetPasswordSendCode(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	pageData := resetpasswordsendcodepage.ResetPasswordSendCodePageData{
		PageData: basePageData,
	}
	app.render(r.Context(), w, r, resetpasswordsendcodepage.Page(pageData))
}

func (app *application) userResetPasswordSendCodePost(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	pageData := resetpasswordsendcodepage.ResetPasswordSendCodePageData{
		PageData: basePageData,
	}
	pageData.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &pageData)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	pageData.CheckFieldError(pageData.Validate.Var(pageData.Email, "required"), "email", "Email field is required")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.Email, "email"), "email", "Email field must be a valid email")

	if !pageData.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, resetpasswordsendcodepage.Page(pageData))
		return
	}

	_, err = app.userService.SendPasswordResetCode(pageData.Email)
	if err != nil && !errors.Is(err, models.ErrUserDoesNotExist) {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "reset-password-email", pageData.Email)
	http.Redirect(w, r, "/user/reset-password/enter-code", http.StatusSeeOther)
}

func (app *application) userResetPasswordEnterCode(w http.ResponseWriter, r *http.Request) {
	email := app.sessionManager.PopString(r.Context(), "reset-password-email")
	basePageData := app.newPageData(r)
	pageData := resetpasswordentercodepage.ResetPasswordEnterCodePageData{
		PageData: basePageData,
		Email:    email,
	}
	app.render(r.Context(), w, r, resetpasswordentercodepage.Page(pageData))
}

func (app *application) userResetPasswordEnterCodePost(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	pageData := resetpasswordentercodepage.ResetPasswordEnterCodePageData{
		PageData: basePageData,
	}
	pageData.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &pageData)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	pageData.CheckFieldError(pageData.Validate.Var(pageData.Code, "required"), "code", "Code field is required")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.Email, "required"), "email", "Email field is required")

	if !pageData.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, resetpasswordentercodepage.Page(pageData))
		return
	}

	valid, err := app.userService.VerifyPasswordResetCode(pageData.Code)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !valid {
		pageData.AddFieldError("code", "Provided code is invalid or has expired")
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, resetpasswordentercodepage.Page(pageData))
		return
	}

	app.sessionManager.Put(r.Context(), "password-reset-code", pageData.Code)
	http.Redirect(w, r, "/user/reset-password", http.StatusSeeOther)
}

func (app *application) userResetPassword(w http.ResponseWriter, r *http.Request) {
	code := app.sessionManager.PopString(r.Context(), "password-reset-code")
	basePageData := app.newPageData(r)
	pageData := resetpasswordpage.ResetPasswordPageData{
		PageData: basePageData,
		Code:     code,
	}

	app.render(r.Context(), w, r, resetpasswordpage.Page(pageData))
}

func (app *application) userResetPasswordPost(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	pageData := resetpasswordpage.ResetPasswordPageData{
		PageData: basePageData,
	}
	pageData.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &pageData)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	pageData.CheckFieldError(pageData.Validate.Var(pageData.Code, "required"), "code", "Code field is required")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.NewPassword, "required"), "new-password", "New password field is required")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.NewPassword, "no_whitespace"), "new-password", "Password cannot contain whitespace characters")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.NewPassword, "min=8"), "new-password", "Must contain at least 8 characters")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.NewPassword, "max=64"), "new-password", "Cannot contain more than 64 characters")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.NewPasswordConfirmation, "required"), "new-password-confirmation", "Field is required")
	pageData.CheckFieldBool(pageData.NewPassword == pageData.NewPasswordConfirmation, "new-password", "Passwords do not match")

	if !pageData.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, resetpasswordpage.Page(pageData))
		return
	}

	request, err := app.userService.ResetPassword(pageData.NewPassword, pageData.Code)

	if err != nil && errors.Is(err, services.ErrPasswordResetCodeInvalid) {
		pageData.AddFieldError("code", "Provided password reset code is invalid")
		w.WriteHeader(http.StatusUnauthorized)
		app.render(r.Context(), w, r, resetpasswordpage.Page(pageData))
		return
	}

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.sessionManager.Iterate(r.Context(), func(ctx context.Context) error {
		userID := app.sessionManager.GetInt32(ctx, "authenticatedUserId")
		if userID == request.UserID {
			return app.sessionManager.Destroy(ctx)
		}
		return nil
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userSettings(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	pageData := usersettingspage.UserSettingsPageData{
		PageData: basePageData,
	}
	pageData.Validate = validator.NewValidator()
	app.render(r.Context(), w, r, usersettingspage.Page(pageData))
}

func (app *application) userSettingsPost(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	pageData := usersettingspage.UserSettingsPageData{
		PageData: basePageData,
	}
	pageData.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &pageData)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if pageData.Form == "change-password" {
		app.handleUserSettingsChangePasswordPost(w, r, &pageData)
		return
	}

	if pageData.Form == "change-email" {
		app.handleUserSettingsChangeEmailPost(w, r, &pageData)
		return
	}

	if pageData.Form == "change-email-send-code" {
		app.handleUserSettingsChangeEmailSendCodePost(w, r, &pageData)
		return
	}

	if pageData.Form == "delete-account" {
		app.handleUserSettingsDeleteAccountPost(w, r, &pageData)
		return
	}

	pageData.AddFieldError("form", "Invalid form field")
	w.WriteHeader(http.StatusBadRequest)
	app.render(r.Context(), w, r, usersettingspage.Page(pageData))
}

func (app *application) handleUserSettingsChangePasswordPost(w http.ResponseWriter, r *http.Request, pageData *usersettingspage.UserSettingsPageData) {
	pageData.CheckFieldError(pageData.Validate.Var(pageData.ChpCurrentPassword, "required"), "chp-current-password", "field is required")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.ChpCurrentPassword, "max=64"), "chp-current-password", "cannot contain more than 64 characters")

	pageData.CheckFieldError(pageData.Validate.Var(pageData.ChpNewPassword, "required"), "chp-new-password", "field is required")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.ChpNewPassword, "no_whitespace"), "chp-new-password", "cannot contain whitespace characters")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.ChpNewPassword, "min=8"), "chp-new-password", "must contain at least 8 characters")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.ChpNewPassword, "max=64"), "chp-new-password", "cannot contain more than 64 characters")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.ChpNewPasswordConfirmation, "required"), "chp-new-password-confirmation", "Field is required")
	pageData.CheckFieldBool(pageData.ChpNewPassword == pageData.ChpNewPasswordConfirmation, "chp-new-password", "passwords do not match")

	if !pageData.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, usersettingspage.Page(*pageData))
		return
	}

	err := app.userService.ChangePassword(pageData.UserInfo.ID, pageData.ChpCurrentPassword, pageData.ChpNewPassword)
	if err != nil && errors.Is(err, services.ErrInvalidCredentials) {
		pageData.AddFieldError("chp-current-password", "incorrect password")
		w.WriteHeader(http.StatusUnauthorized)
		app.render(r.Context(), w, r, usersettingspage.Page(*pageData))
		return
	}

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.sessionManager.Iterate(r.Context(), func(ctx context.Context) error {
		userID := app.sessionManager.GetInt32(ctx, "authenticatedUserId")
		if userID == pageData.UserInfo.ID {
			return app.sessionManager.Destroy(ctx)
		}
		return nil
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) handleUserSettingsChangeEmailPost(w http.ResponseWriter, r *http.Request, pageData *usersettingspage.UserSettingsPageData) {
	pageData.CheckFieldError(pageData.Validate.Var(pageData.ChePassword, "required"), "che-password", "field is required")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.ChePassword, "max=64"), "che-password", "cannot contain more than 64 characters")

	pageData.CheckFieldError(pageData.Validate.Var(pageData.CheNewEmail, "required"), "che-new-email", "field is required")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.CheNewEmail, "email"), "che-new-email", "field must be a valid email")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.CheNewEmail, "max=128"), "che-new-email", "cannot contain more than 128 characters")
	pageData.CheckFieldBool(pageData.CheNewEmail != pageData.UserInfo.Email, "che-new-email", "cannot be your current email")

	if !pageData.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, usersettingspage.Page(*pageData))
		return
	}

	err := app.userService.SendEmailChangeRequestCode(pageData.UserInfo.ID, pageData.ChePassword, pageData.CheNewEmail)
	if err != nil && errors.Is(err, services.ErrInvalidCredentials) {
		pageData.AddFieldError("che-password", "incorrect password")
		w.WriteHeader(http.StatusUnauthorized)
		app.render(r.Context(), w, r, usersettingspage.Page(*pageData))
		return
	}

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	pageData.Form = "change-email-send-code"
	app.render(r.Context(), w, r, usersettingspage.Page(*pageData))
}

func (app *application) handleUserSettingsChangeEmailSendCodePost(w http.ResponseWriter, r *http.Request, pageData *usersettingspage.UserSettingsPageData) {
	pageData.CheckFieldError(pageData.Validate.Var(pageData.ChescCode, "required"), "chesc-code", "field is required")

	if !pageData.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, usersettingspage.Page(*pageData))
		return
	}

	ok, err := app.userService.ChangeEmail(pageData.ChescCode)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !ok {
		pageData.AddFieldError("chesc-code", "code incorrect or expired")
		w.WriteHeader(http.StatusUnauthorized)
		app.render(r.Context(), w, r, usersettingspage.Page(*pageData))
		return
	}

	http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
}

func (app *application) handleUserSettingsDeleteAccountPost(w http.ResponseWriter, r *http.Request, pageData *usersettingspage.UserSettingsPageData) {
	pageData.CheckFieldError(pageData.Validate.Var(pageData.DaEmail, "required"), "da-email", "field is required")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.DaEmail, "email"), "da-email", "field must be a valid email")
	pageData.CheckFieldError(pageData.Validate.Var(pageData.DaEmail, "max=128"), "da-email", "cannot contain more than 128 characters")
	pageData.CheckFieldBool(pageData.DaEmail == pageData.UserInfo.Email, "da-email", "type in your current email")

	if !pageData.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, usersettingspage.Page(*pageData))
		return
	}

	err := app.userService.DeleteAccount(pageData.UserInfo.ID)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
