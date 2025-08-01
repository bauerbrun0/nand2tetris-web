package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/internal/validator"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/resetpasswordentercodepage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/resetpasswordpage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/resetpasswordsendcodepage"
)

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
