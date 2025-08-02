package handlers

import (
	"errors"
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/usersettingspage"
)

func (h *Handlers) handleUserSettingsChangeEmailPost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	data.CheckFieldTag(data.ChangeEmail.Password, "required", "ChangeEmail.Password", data.T("error.field_required"))
	data.CheckFieldTag(
		data.ChangeEmail.Password, "max=64", "ChangeEmail.Password", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}),
	)

	data.CheckFieldTag(data.ChangeEmail.NewEmail, "required", "ChangeEmail.NewEmail", data.T("error.field_required"))
	data.CheckFieldTag(data.ChangeEmail.NewEmail, "email", "ChangeEmail.NewEmail", data.T("error.field_invalid_email"))
	data.CheckFieldTag(
		data.ChangeEmail.NewEmail, "max=128", "ChangeEmail.NewEmail", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "128"}),
	)
	data.CheckFieldBool(
		data.ChangeEmail.NewEmail != data.UserInfo.Email, "ChangeEmail.NewEmail", data.T("error.field_cannot_be_current_email"),
	)

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	err := h.UserService.SendEmailChangeRequestCode(data.UserInfo.ID, data.ChangeEmail.Password, data.ChangeEmail.NewEmail)
	if err != nil && errors.Is(err, services.ErrInvalidCredentials) {
		data.AddFieldError("ChangeEmail.Password", data.T("error.incorrect_password"))
		w.WriteHeader(http.StatusUnauthorized)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	data.Action = usersettingspage.Action(ActionChangeEmailSendCode)
	h.Render(r.Context(), w, r, usersettingspage.Page(*data))
}

func (h *Handlers) handleUserSettingsChangeEmailSendCodePost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	data.CheckFieldTag(data.ChangeEmailSendCode.Code, "required", "ChangeEmailSendCode.Code", data.T("error.field_required"))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	ok, err := h.UserService.ChangeEmail(data.ChangeEmailSendCode.Code)
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	if !ok {
		data.AddFieldError("ChangeEmailSendCode.Code", data.T("error.code_incorrect"))
		w.WriteHeader(http.StatusUnauthorized)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.T("toast.successfully_changed_email"),
			Variant:  "success",
			Duration: 2000,
		},
	})
	http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
}
