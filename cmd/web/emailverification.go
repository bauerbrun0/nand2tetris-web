package main

import (
	"errors"
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/internal/validator"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/verifyemailpage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/verifyemailsendcodepage"
)

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
