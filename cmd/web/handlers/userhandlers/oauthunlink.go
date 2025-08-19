package userhandlers

import (
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/usersettingspage"
)

func (h *Handlers) unlinkOAuthAccount(w http.ResponseWriter, r *http.Request, data *pages.PageData, provider models.Provider) {
	if len(data.UserInfo.LinkedAccounts) == 1 && !data.UserInfo.IsPasswordSet {
		h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.only_login_method"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
		return
	}

	err := h.UserService.RemoveOAuthAuthorization(data.UserInfo.ID, provider)
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.TTemplate("toast.oauth_user_successfully_unlinked", map[string]string{"Provider": string(provider)}),
			Variant:  "success",
			Duration: 3000,
		},
	})
	http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
}

func (h *Handlers) handleUserSettingsUnlinkGoogleAccountPost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	verificationMethod, ok := ParseVerificationMethod(data.Verification)
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	switch verificationMethod {
	case VerificationPassword:
		ok := h.validateAndCheckPasswordField(w, r, data, data.UnlinkGoogle.Password, "UnlinkGoogle.Password")
		if ok {
			newPageData := h.NewPageData(r)
			h.unlinkOAuthAccount(w, r, &newPageData, models.ProviderGoogle)
		}
	case VerificationGoogle:
		h.sendGoogleActionRedirect(w, r, ActionUnlinkGoogleAccount, "/user/oauth/google/callback/action")
	case VerificationGitHub:
		h.sendGithubActionRedirect(w, r, ActionUnlinkGoogleAccount, "/user/oauth/github/callback/action")
	default:
		w.WriteHeader(http.StatusBadRequest)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
	}
}

func (h *Handlers) handleUserSettingsUnlinkGithubAccountPost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	verificationMethod, ok := ParseVerificationMethod(data.Verification)
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	switch verificationMethod {
	case VerificationPassword:
		ok := h.validateAndCheckPasswordField(w, r, data, data.UnlinkGithub.Password, "UnlinkGithub.Password")
		if ok {
			newPageData := h.NewPageData(r)
			h.unlinkOAuthAccount(w, r, &newPageData, models.ProviderGitHub)
		}
	case VerificationGoogle:
		h.sendGoogleActionRedirect(w, r, ActionUnlinkGitHubAccount, "/user/oauth/google/callback/action")
	case VerificationGitHub:
		h.sendGithubActionRedirect(w, r, ActionUnlinkGitHubAccount, "/user/oauth/github/callback/action")
	default:
		w.WriteHeader(http.StatusBadRequest)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
	}
}
