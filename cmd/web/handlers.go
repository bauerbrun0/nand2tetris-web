package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/internal/crypto"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/internal/validator"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
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
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newPageData(r)
	if app.isAuthenticated(r) {
		app.render(r.Context(), w, r, homepage.Page(data))
		return
	}
	app.render(r.Context(), w, r, landingpage.Page(data))
}

func (app *application) userRegister(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	data := registerpage.RegisterPageData{
		PageData: basePageData,
	}
	app.render(r.Context(), w, r, registerpage.Page(data))
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	data := loginpage.LoginPageData{
		PageData: basePageData,
	}
	app.render(r.Context(), w, r, loginpage.Page(data))
}

func (app *application) userRegisterPost(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	data := registerpage.RegisterPageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	data.CheckFieldTag(data.Email, "required", "email", data.T("error.field_required"))
	data.CheckFieldTag(data.Email, "email", "email", data.T("error.field_invalid_email"))
	data.CheckFieldTag(data.Email, "max=128", "email", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "128"}))
	data.CheckFieldTag(data.Username, "required", "username", data.T("error.field_required"))
	data.CheckFieldTag(data.Username, "min=3", "username", data.TTemplate("error.field_not_enough_characters", map[string]string{"Min": "3"}))
	data.CheckFieldTag(data.Username, "max=64", "username", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}))
	err = data.Validate.Var(data.Username, "email")
	if err == nil {
		data.AddFieldError("username", data.T("error.field_is_email"))
	}
	data.CheckFieldTag(data.Username, "no_whitespace", "username", data.T("error.field_contains_whitespace"))
	data.CheckFieldTag(data.Password, "required", "password", data.T("error.field_required"))
	data.CheckFieldTag(data.Password, "no_whitespace", "password", data.T("error.field_contains_whitespace"))
	data.CheckFieldTag(data.Password, "min=8", "password", data.TTemplate("error.field_not_enough_characters", map[string]string{"Min": "8"}))
	data.CheckFieldTag(data.Password, "max=64", "password", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}))
	data.CheckFieldTag(data.PasswordConfirmation, "required", "password-confirmation", data.T("error.field_required"))
	data.CheckFieldBool(data.Password == data.PasswordConfirmation, "password", data.T("error.passwords_do_not_match"))
	data.CheckFieldTag(data.Terms, "required,eq=on", "terms", data.T("error.terms_conditions"))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, registerpage.Page(data))
		return
	}

	_, err = app.userService.CreateUser(data.Email, data.Username, data.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			data.AddFieldError("email", data.T("error.email_already_used"))
			w.WriteHeader(http.StatusUnprocessableEntity)
			app.render(r.Context(), w, r, registerpage.Page(data))
			return
		}
		if errors.Is(err, models.ErrDuplicateUsername) {
			data.AddFieldError("username", data.T("error.username_already_used"))
			w.WriteHeader(http.StatusUnprocessableEntity)
			app.render(r.Context(), w, r, registerpage.Page(data))
			return
		}
		app.serverError(w, r, err)
	}

	app.sessionManager.Put(r.Context(), "email-to-verify", data.Email)
	http.Redirect(w, r, "/user/verify-email", http.StatusSeeOther)
}

func (app *application) userVerifyEmail(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	email := app.sessionManager.PopString(r.Context(), "email-to-verify")
	data := verifyemailpage.VerifyEmailPageData{
		PageData: basePageData,
		Email:    email,
	}
	app.render(r.Context(), w, r, verifyemailpage.Page(data))
}

func (app *application) userVerifyEmailPost(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	data := verifyemailpage.VerifyEmailPageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	data.CheckFieldTag(data.Code, "required", "code", data.TTemplate("error.x_field_required", map[string]string{"Field": "Code"}))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, verifyemailpage.Page(data))
		return
	}

	ok, err := app.userService.VerifyEmail(data.Code)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !ok {
		data.AddFieldError("code", data.T("error.verification_code_invalid"))
		app.render(r.Context(), w, r, verifyemailpage.Page(data))
		return
	}
	app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.T("toast.successfully_registered"),
			Variant:  "success",
			Duration: 3000,
		},
	})
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userVerifyEmailResendCode(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	data := verifyemailsendcodepage.VerifyEmailSendCodePageData{
		PageData: basePageData,
	}
	app.render(r.Context(), w, r, verifyemailsendcodepage.Page(data))
}

func (app *application) userVerifyEmailResendCodePost(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	data := verifyemailsendcodepage.VerifyEmailSendCodePageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	data.CheckFieldTag(data.Email, "required", "email", data.TTemplate("error.x_field_required", map[string]string{"Field": "Email"}))
	data.CheckFieldTag(data.Email, "email", "email", data.T("error.email_field_invalid_email"))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, verifyemailsendcodepage.Page(data))
		return
	}

	_, err = app.userService.ResendEmailVerificationCode(data.Email)

	if err != nil && !errors.Is(err, models.ErrUserDoesNotExist) && !errors.Is(err, services.ErrEmailAlreadyVerified) {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "email-to-verify", data.Email)
	app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.TTemplate("toast.new_email_sent_to", map[string]string{"Email": data.Email}),
			Variant:  "info",
			Duration: 3000,
		},
	})
	http.Redirect(w, r, "/user/verify-email", http.StatusSeeOther)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	data := loginpage.LoginPageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	data.CheckFieldTag(data.Username, "required", "username", data.T("error.field_required"))
	data.CheckFieldTag(data.Password, "required", "password", data.T("error.field_required"))
	data.CheckFieldTag(data.Password, "max=64", "password", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, loginpage.Page(data))
		return
	}

	user, err := app.userService.AuthenticateUser(data.Username, data.Password)
	if err != nil && errors.Is(err, services.ErrInvalidCredentials) {
		data.AddFieldError("password", data.T("error.invalid_credentials"))
		w.WriteHeader(http.StatusUnauthorized)
		app.render(r.Context(), w, r, loginpage.Page(data))
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
	app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.TTemplate("toast.user_welcome", map[string]string{"Username": user.Username}),
			Variant:  "simple",
			Duration: 2000,
		},
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserId")
	app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  app.getLocalizer(r).MustLocalize(&i18n.LocalizeConfig{MessageID: "toast.logout"}),
			Variant:  "simple",
			Duration: 2000,
		},
	})
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userResetPasswordSendCode(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	data := resetpasswordsendcodepage.ResetPasswordSendCodePageData{
		PageData: basePageData,
	}
	app.render(r.Context(), w, r, resetpasswordsendcodepage.Page(data))
}

func (app *application) userResetPasswordSendCodePost(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	data := resetpasswordsendcodepage.ResetPasswordSendCodePageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data.CheckFieldTag(data.Email, "required", "email", data.TTemplate("error.x_field_required", map[string]string{"Field": "Email"}))
	data.CheckFieldTag(data.Email, "email", "email", data.T("error.email_field_invalid_email"))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, resetpasswordsendcodepage.Page(data))
		return
	}

	_, err = app.userService.SendPasswordResetCode(data.Email)
	if err != nil && !errors.Is(err, models.ErrUserDoesNotExist) {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "reset-password-email", data.Email)
	app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.TTemplate("toast.email_sent_to", map[string]string{"Email": data.Email}),
			Variant:  "info",
			Duration: 3000,
		},
	})
	http.Redirect(w, r, "/user/reset-password/enter-code", http.StatusSeeOther)
}

func (app *application) userResetPasswordEnterCode(w http.ResponseWriter, r *http.Request) {
	email := app.sessionManager.PopString(r.Context(), "reset-password-email")
	basePageData := app.newPageData(r)
	data := resetpasswordentercodepage.ResetPasswordEnterCodePageData{
		PageData: basePageData,
		Email:    email,
	}
	app.render(r.Context(), w, r, resetpasswordentercodepage.Page(data))
}

func (app *application) userResetPasswordEnterCodePost(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	data := resetpasswordentercodepage.ResetPasswordEnterCodePageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data.CheckFieldTag(data.Code, "required", "code", data.TTemplate("error.x_field_required", map[string]string{"Field": "Code"}))
	data.CheckFieldTag(data.Email, "required", "email", data.TTemplate("error.x_field_required", map[string]string{"Field": "Email"}))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, resetpasswordentercodepage.Page(data))
		return
	}

	valid, err := app.userService.VerifyPasswordResetCode(data.Code)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !valid {
		data.AddFieldError("code", data.T("error.provided_code_invalid"))
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, resetpasswordentercodepage.Page(data))
		return
	}

	app.sessionManager.Put(r.Context(), "password-reset-code", data.Code)
	http.Redirect(w, r, "/user/reset-password", http.StatusSeeOther)
}

func (app *application) userResetPassword(w http.ResponseWriter, r *http.Request) {
	code := app.sessionManager.PopString(r.Context(), "password-reset-code")
	basePageData := app.newPageData(r)
	data := resetpasswordpage.ResetPasswordPageData{
		PageData: basePageData,
		Code:     code,
	}

	app.render(r.Context(), w, r, resetpasswordpage.Page(data))
}

func (app *application) userResetPasswordPost(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	data := resetpasswordpage.ResetPasswordPageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data.CheckFieldTag(data.Code, "required", "code", data.TTemplate("error.x_field_required", map[string]string{"Field": "Code"}))
	data.CheckFieldTag(
		data.NewPassword, "required", "new-password", data.TTemplate("error.x_field_required", map[string]string{"Field": "New password"}),
	)
	data.CheckFieldTag(data.NewPassword, "no_whitespace", "new-password", data.T("error.password_contains_whitespace"))
	data.CheckFieldTag(data.NewPassword, "min=8", "new-password", data.TTemplate("error.field_not_enough_characters", map[string]string{"Min": "8"}))
	data.CheckFieldTag(data.NewPassword, "max=64", "new-password", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}))
	data.CheckFieldTag(data.NewPasswordConfirmation, "required", "new-password-confirmation", data.T("error.field_required"))
	data.CheckFieldBool(data.NewPassword == data.NewPasswordConfirmation, "new-password", data.T("error.passwords_do_not_match"))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, resetpasswordpage.Page(data))
		return
	}

	request, err := app.userService.ResetPassword(data.NewPassword, data.Code)

	if err != nil && errors.Is(err, services.ErrPasswordResetCodeInvalid) {
		data.AddFieldError("code", data.T("error.provided_password_reset_code_invalid"))
		w.WriteHeader(http.StatusUnauthorized)
		app.render(r.Context(), w, r, resetpasswordpage.Page(data))
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

	app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.T("toast.successfully_changed_password"),
			Variant:  "success",
			Duration: 3000,
		},
	})
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userSettings(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	data := usersettingspage.UserSettingsPageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()
	app.render(r.Context(), w, r, usersettingspage.Page(data))
}

func (app *application) userSettingsPost(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	data := usersettingspage.UserSettingsPageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if data.Form == "change-password" {
		app.handleUserSettingsChangePasswordPost(w, r, &data)
		return
	}

	if data.Form == "change-email" {
		app.handleUserSettingsChangeEmailPost(w, r, &data)
		return
	}

	if data.Form == "change-email-send-code" {
		app.handleUserSettingsChangeEmailSendCodePost(w, r, &data)
		return
	}

	if data.Form == "delete-account" {
		app.handleUserSettingsDeleteAccountPost(w, r, &data)
		return
	}

	if data.Form == "create-password" {
		app.handleUserSettingsCreatePasswordPost(w, r, &data)
		return
	}

	data.AddFieldError("form", data.T("error.invalid_form_field"))
	w.WriteHeader(http.StatusBadRequest)
	app.render(r.Context(), w, r, usersettingspage.Page(data))
}

func (app *application) handleUserSettingsChangePasswordPost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	data.CheckFieldTag(data.ChpCurrentPassword, "required", "chp-current-password", data.T("error.field_required"))
	data.CheckFieldTag(
		data.ChpCurrentPassword, "max=64", "chp-current-password", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}),
	)

	data.CheckFieldTag(data.ChpNewPassword, "required", "chp-new-password", data.T("error.field_required"))
	data.CheckFieldTag(data.ChpNewPassword, "no_whitespace", "chp-new-password", data.T("error.field_contains_whitespace"))
	data.CheckFieldTag(
		data.ChpNewPassword, "min=8", "chp-new-password", data.TTemplate("error.field_not_enough_characters", map[string]string{"Min": "8"}),
	)
	data.CheckFieldTag(
		data.ChpNewPassword, "max=64", "chp-new-password", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}),
	)
	data.CheckFieldTag(data.ChpNewPasswordConfirmation, "required", "chp-new-password-confirmation", data.T("error.field_required"))
	data.CheckFieldBool(
		data.ChpNewPassword == data.ChpNewPasswordConfirmation, "chp-new-password", data.T("error.passwords_do_not_match"),
	)

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	err := app.userService.ChangePassword(data.UserInfo.ID, data.ChpCurrentPassword, data.ChpNewPassword)
	if err != nil && errors.Is(err, services.ErrInvalidCredentials) {
		data.AddFieldError("chp-current-password", data.T("error.incorrect_password"))
		w.WriteHeader(http.StatusUnauthorized)
		app.render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.sessionManager.Iterate(r.Context(), func(ctx context.Context) error {
		userID := app.sessionManager.GetInt32(ctx, "authenticatedUserId")
		if userID == data.UserInfo.ID {
			return app.sessionManager.Destroy(ctx)
		}
		return nil
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Destroy(r.Context())
	app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.T("toast.successfully_changed_password"),
			Variant:  "success",
			Duration: 3000,
		},
	})

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) handleUserSettingsChangeEmailPost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	data.CheckFieldTag(data.ChePassword, "required", "che-password", data.T("error.field_required"))
	data.CheckFieldTag(
		data.ChePassword, "max=64", "che-password", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}),
	)

	data.CheckFieldTag(data.CheNewEmail, "required", "che-new-email", data.T("error.field_required"))
	data.CheckFieldTag(data.CheNewEmail, "email", "che-new-email", data.T("error.field_invalid_email"))
	data.CheckFieldTag(
		data.CheNewEmail, "max=128", "che-new-email", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "128"}),
	)
	data.CheckFieldBool(
		data.CheNewEmail != data.UserInfo.Email, "che-new-email", data.T("error.field_cannot_be_current_email"),
	)

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	err := app.userService.SendEmailChangeRequestCode(data.UserInfo.ID, data.ChePassword, data.CheNewEmail)
	if err != nil && errors.Is(err, services.ErrInvalidCredentials) {
		data.AddFieldError("che-password", data.T("error.incorrect_password"))
		w.WriteHeader(http.StatusUnauthorized)
		app.render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data.Form = "change-email-send-code"
	app.render(r.Context(), w, r, usersettingspage.Page(*data))
}

func (app *application) handleUserSettingsChangeEmailSendCodePost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	data.CheckFieldTag(data.ChescCode, "required", "chesc-code", data.T("error.field_required"))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	ok, err := app.userService.ChangeEmail(data.ChescCode)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !ok {
		data.AddFieldError("chesc-code", data.T("error.code_incorrect"))
		w.WriteHeader(http.StatusUnauthorized)
		app.render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.T("toast.successfully_changed_email"),
			Variant:  "success",
			Duration: 2000,
		},
	})
	http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
}

func (app *application) handleUserSettingsDeleteAccountPost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	data.CheckFieldTag(data.DaEmail, "required", "da-email", data.T("error.field_required"))
	data.CheckFieldTag(data.DaEmail, "email", "da-email", data.T("error.field_invalid_email"))
	data.CheckFieldTag(data.DaEmail, "max=128", "da-email", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "128"}))
	data.CheckFieldBool(data.DaEmail == data.UserInfo.Email, "da-email", data.T("error.type_current_email"))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	err := app.userService.DeleteAccount(data.UserInfo.ID)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.T("toast.successfully_deleted_account"),
			Variant:  "success",
			Duration: 2000,
		},
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) handleUserSettingsCreatePasswordPost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	data.CheckFieldTag(data.CpaPassword, "required", "cpa-password", data.T("error.field_required"))
	data.CheckFieldTag(data.CpaPassword, "no_whitespace", "cpa-password", data.T("error.field_contains_whitespace"))
	data.CheckFieldTag(
		data.CpaPassword, "min=8", "cpa-password", data.TTemplate("error.field_not_enough_characters", map[string]string{"Min": "8"}),
	)
	data.CheckFieldTag(
		data.CpaPassword, "max=64", "cpa-password", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}),
	)
	data.CheckFieldTag(data.CpaPasswordConfirmation, "required", "cpa-password-confirmation", data.T("error.field_required"))
	data.CheckFieldBool(
		data.CpaPassword == data.CpaPasswordConfirmation, "cpa-password", data.T("error.passwords_do_not_match"),
	)

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	err := app.userService.CreatePassword(data.UserInfo.ID, data.CpaPassword)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.sessionManager.Iterate(r.Context(), func(ctx context.Context) error {
		userID := app.sessionManager.GetInt32(ctx, "authenticatedUserId")
		if userID == data.UserInfo.ID {
			return app.sessionManager.Destroy(ctx)
		}
		return nil
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Destroy(r.Context())
	app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.T("toast.successfully_created_password"),
			Variant:  "success",
			Duration: 3000,
		},
	})

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLoginGithub(w http.ResponseWriter, r *http.Request) {
	state := crypto.GenerateRandomString(16)
	c := &http.Cookie{
		Name:     "github_state",
		Value:    state,
		Path:     "/",
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	redirectUrl := app.githubOauthService.GetRedirectUrl(state)
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (app *application) userLoginGithubCallback(w http.ResponseWriter, r *http.Request) {
	data := app.newPageData(r)

	state, err := r.Cookie("github_state")
	if err != nil {
		app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.state_cookie_not_found"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/login", http.StatusBadRequest)
		return
	}

	queryState := r.URL.Query().Get("state")
	if state.Value != queryState {
		app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.state_tokens_do_not_match"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/login", http.StatusUnauthorized)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := app.githubOauthService.ExchangeCodeForToken(code)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	oauthUser, err := app.githubOauthService.GetUserInfo(token)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	user, err := app.userService.AuthenticateOAuthUser(oauthUser, models.ProviderGitHub)
	if err != nil && !errors.Is(err, services.ErrUserAlreadyExists) {
		app.serverError(w, r, err)
		return
	}

	if err != nil && errors.Is(err, services.ErrUserAlreadyExists) {
		app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.oauth_user_exists"),
				Variant:  "error",
				Duration: 5000,
			},
		})
		http.Redirect(w, r, "/user/login", http.StatusUnauthorized)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserId", user.ID)
	app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.TTemplate("toast.user_welcome", map[string]string{"Username": user.Username}),
			Variant:  "simple",
			Duration: 2000,
		},
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLoginGoogle(w http.ResponseWriter, r *http.Request) {
	state := crypto.GenerateRandomString(30)
	c := &http.Cookie{
		Name:     "google_state",
		Value:    state,
		Path:     "/",
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	redirectUrl := app.googleOauthService.GetRedirectUrl(state)
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (app *application) userLoginGoogleCallback(w http.ResponseWriter, r *http.Request) {
	data := app.newPageData(r)

	state, err := r.Cookie("google_state")
	if err != nil {
		app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.state_cookie_not_found"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	queryState := r.URL.Query().Get("state")
	if state.Value != queryState {
		app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.state_tokens_do_not_match"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := app.googleOauthService.ExchangeCodeForToken(code)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	oauthUser, err := app.googleOauthService.GetUserInfo(token)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	user, err := app.userService.AuthenticateOAuthUser(oauthUser, models.ProviderGoogle)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserId", user.ID)
	app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.TTemplate("toast.user_welcome", map[string]string{"Username": user.Username}),
			Variant:  "simple",
			Duration: 2000,
		},
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
