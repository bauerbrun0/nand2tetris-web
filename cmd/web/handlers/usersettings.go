package handlers

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

	switch data.Form {
	case "change-password":
		h.handleUserSettingsChangePasswordPost(w, r, &data)
	case "change-email":
		h.handleUserSettingsChangeEmailPost(w, r, &data)
	case "change-email-send-code":
		h.handleUserSettingsChangeEmailSendCodePost(w, r, &data)
	case "delete-account":
		h.handleUserSettingsDeleteAccountPost(w, r, &data)
	case "create-password":
		h.handleUserSettingsCreatePasswordPost(w, r, &data)
	case "link-google-account":
		h.handleUserSettingsLinkGoogleAccountPost(w, r, &data)
	case "link-github-account":
		h.handleUserSettingsLinkGithubAccountPost(w, r, &data)
	case "unlink-google-account":
		h.handleUserSettingsUnlinkGoogleAccountPost(w, r, &data)
	case "unlink-github-account":
		h.handleUserSettingsUnlinkGithubAccountPost(w, r, &data)
	default:
		w.WriteHeader(http.StatusBadRequest)
		h.Render(r.Context(), w, r, usersettingspage.Page(data))
	}
}
