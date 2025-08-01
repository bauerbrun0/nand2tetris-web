package main

import (
	"errors"
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/usersettingspage"
)

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
