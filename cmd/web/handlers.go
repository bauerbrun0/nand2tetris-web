package main

import (
	"errors"
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/validator"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/landingpage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/registerpage"
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

	http.Redirect(w, r, "/user/verify-email", http.StatusSeeOther)
}
