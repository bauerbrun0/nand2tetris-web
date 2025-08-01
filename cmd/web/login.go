package main

import (
	"errors"
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/internal/validator"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/loginpage"
)

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	data := loginpage.LoginPageData{
		PageData: basePageData,
	}
	app.render(r.Context(), w, r, loginpage.Page(data))
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
