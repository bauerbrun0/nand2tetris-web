package userhandlers

import (
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/validator"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/usersettingspage"
)

func (h *Handlers) UserSettings(w http.ResponseWriter, r *http.Request) {
	basePageData := h.NewPageData(r)
	data := usersettingspage.UserSettingsPageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()
	h.Render(r.Context(), w, r, usersettingspage.Page(data))
}

func (h *Handlers) UserSettingsPost(w http.ResponseWriter, r *http.Request) {
	basePageData := h.NewPageData(r)
	data := usersettingspage.UserSettingsPageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()

	err := h.DecodePostForm(r, &data)
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	action, ok := ParseAction(data.Action)
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, usersettingspage.Page(data))
		return
	}

	switch action {
	case ActionChangePassword:
		h.handleUserSettingsChangePasswordPost(w, r, &data)
	case ActionChangeEmail:
		h.handleUserSettingsChangeEmailPost(w, r, &data)
	case ActionChangeEmailSendCode:
		h.handleUserSettingsChangeEmailSendCodePost(w, r, &data)
	case ActionDeleteAccount:
		h.handleUserSettingsDeleteAccountPost(w, r, &data)
	case ActionCreatePassword:
		h.handleUserSettingsCreatePasswordPost(w, r, &data)
	case ActionLinkGoogleAccount:
		h.handleUserSettingsLinkGoogleAccountPost(w, r, &data)
	case ActionLinkGitHubAccount:
		h.handleUserSettingsLinkGithubAccountPost(w, r, &data)
	case ActionUnlinkGoogleAccount:
		h.handleUserSettingsUnlinkGoogleAccountPost(w, r, &data)
	case ActionUnlinkGitHubAccount:
		h.handleUserSettingsUnlinkGithubAccountPost(w, r, &data)
	default:
		w.WriteHeader(http.StatusBadRequest)
		h.Render(r.Context(), w, r, usersettingspage.Page(data))
	}
}
