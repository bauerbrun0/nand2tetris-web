package main

import (
	"errors"
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/internal/validator"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/landingpage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/loginpage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/registerpage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/verifyemailpage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/verifyemailsendcodepage"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	msg := app.sessionManager.GetString(r.Context(), "message")
	app.logger.Info("getting message from session", "message", msg)
	app.sessionManager.Put(r.Context(), "message", "Hello World!")
	app.logger.Info("stored message to session")
	app.render(r.Context(), w, r, landingpage.Page())
}

func (app *application) userRegister(w http.ResponseWriter, r *http.Request) {
	app.render(r.Context(), w, r, registerpage.Page(&registerpage.RegisterForm{}))
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	app.render(r.Context(), w, r, loginpage.Page(&loginpage.LoginForm{}))
}

func (app *application) userRegisterPost(w http.ResponseWriter, r *http.Request) {
	var form registerpage.RegisterForm
	form.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &form)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	form.CheckFieldError(form.Validate.Var(form.Email, "required"), "email", "field is required")
	form.CheckFieldError(form.Validate.Var(form.Email, "email"), "email", "field must be a valid email")
	form.CheckFieldError(form.Validate.Var(form.Username, "required"), "username", "field is required")
	form.CheckFieldError(form.Validate.Var(form.Password, "required"), "password", "field is required")
	form.CheckFieldError(form.Validate.Var(form.Password, "min=8"), "password", "must contain at least 8 characters")
	form.CheckFieldError(form.Validate.Var(form.PasswordConfirmation, "required"), "password-confirmation", "field is required")
	form.CheckFieldBool(form.Password == form.PasswordConfirmation, "password", "passwords do not match")
	form.CheckFieldError(form.Validate.Var(form.Terms, "required,eq=on"), "terms", "You must agree to the Terms & Conditions")

	if !form.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, registerpage.Page(&form))
		return
	}

	_, err = app.userService.CreateUser(form.Email, form.Username, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			w.WriteHeader(http.StatusUnprocessableEntity)
			app.render(r.Context(), w, r, registerpage.Page(&form))
			return
		}
		if errors.Is(err, models.ErrDuplicateUsername) {
			form.AddFieldError("username", "Username is already in use")
			w.WriteHeader(http.StatusUnprocessableEntity)
			app.render(r.Context(), w, r, registerpage.Page(&form))
			return
		}
		app.serverError(w, r, err)
	}

	app.sessionManager.Put(r.Context(), "email-to-verify", form.Email)
	http.Redirect(w, r, "/user/verify-email", http.StatusSeeOther)
}

func (app *application) userVerifyEmail(w http.ResponseWriter, r *http.Request) {
	email := app.sessionManager.PopString(r.Context(), "email-to-verify")
	app.render(r.Context(), w, r, verifyemailpage.Page(email, &verifyemailpage.VerifyEmailForm{}))
}

func (app *application) userVerifyEmailPost(w http.ResponseWriter, r *http.Request) {
	var form verifyemailpage.VerifyEmailForm
	form.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &form)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	form.CheckFieldError(form.Validate.Var(form.Code, "required"), "code", "Code field is required")

	if !form.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, verifyemailpage.Page(form.Email, &form))
		return
	}

	ok, err := app.userService.VerifyEmail(form.Code)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !ok {
		form.AddFieldError("code", "The verification code you entered is invalid or has expired. Please try again.")
		app.render(r.Context(), w, r, verifyemailpage.Page(form.Email, &form))
		return
	}
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userVerifyEmailResendCode(w http.ResponseWriter, r *http.Request) {
	app.render(r.Context(), w, r, verifyemailsendcodepage.Page(&verifyemailsendcodepage.VerifyEmailSendCodeForm{}))
}

func (app *application) userVerifyEmailResendCodePost(w http.ResponseWriter, r *http.Request) {
	var form verifyemailsendcodepage.VerifyEmailSendCodeForm
	form.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &form)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	form.CheckFieldError(form.Validate.Var(form.Email, "required"), "email", "Email field is required")
	form.CheckFieldError(form.Validate.Var(form.Email, "email"), "email", "Email field must be a valid email")

	if !form.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, verifyemailsendcodepage.Page(&form))
		return
	}

	_, err = app.userService.ResendEmailVerificationCode(form.Email)

	if err != nil && !errors.Is(err, models.ErrUserDoesNotExist) && !errors.Is(err, services.ErrEmailAlreadyVerified) {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "email-to-verify", form.Email)
	http.Redirect(w, r, "/user/verify-email", http.StatusSeeOther)
}
