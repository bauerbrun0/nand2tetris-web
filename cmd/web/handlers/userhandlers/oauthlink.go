package userhandlers

import (
	"net/http"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/internal/crypto"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/usersettingspage"
)

func (h *Handlers) UserLinkGithubCallback(w http.ResponseWriter, r *http.Request) {
	data := h.NewPageData(r)

	state, err := r.Cookie("github_state")
	if err != nil {
		h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.state_cookie_not_found"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
		return
	}

	queryState := r.URL.Query().Get("state")
	if state.Value != queryState {
		h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.state_tokens_do_not_match"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := h.GithubOauthService.ExchangeCodeForToken(services.TokenExchangeOptions{
		Code: code,
	})
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	oauthUser, err := h.GithubOauthService.GetUserInfo(token)
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	_, exists, err := h.UserService.GetUserIdByUserProviderId(models.ProviderGitHub, oauthUser.Id)
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	if exists {
		h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.oauth_user_already_linked"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
		return
	}

	err = h.UserService.AddOAuthAuthorization(oauthUser.Id, data.UserInfo.ID, models.ProviderGitHub)
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.TTemplate("toast.oauth_user_successfully_linked", map[string]string{"Provider": "GitHub"}),
			Variant:  "success",
			Duration: 3000,
		},
	})
	http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
}

func (h *Handlers) UserLinkGoogleCallback(w http.ResponseWriter, r *http.Request) {
	data := h.NewPageData(r)

	state, err := r.Cookie("google_state")
	if err != nil {
		h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.state_cookie_not_found"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
		return
	}

	queryState := r.URL.Query().Get("state")
	if state.Value != queryState {
		h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.state_tokens_do_not_match"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := h.GoogleOauthService.ExchangeCodeForToken(services.TokenExchangeOptions{
		Code:         code,
		RedirectPath: "/user/oauth/google/callback/link",
	})
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	oauthUser, err := h.GoogleOauthService.GetUserInfo(token)
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	_, exists, err := h.UserService.GetUserIdByUserProviderId(models.ProviderGoogle, oauthUser.Id)
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	if exists {
		h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.oauth_user_already_linked"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
		return
	}

	err = h.UserService.AddOAuthAuthorization(oauthUser.Id, data.UserInfo.ID, models.ProviderGoogle)
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.TTemplate("toast.oauth_user_successfully_linked", map[string]string{"Provider": "Google"}),
			Variant:  "success",
			Duration: 3000,
		},
	})
	http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
}

func (h *Handlers) handleUserSettingsLinkGoogleAccountPost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	verificationMethod, ok := ParseVerificationMethod(data.Verification)
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	switch verificationMethod {
	case VerificationPassword:
		ok := h.validateAndCheckPasswordField(w, r, data, data.LinkGoogle.Password, "LinkGoogle.Password")
		if ok {
			h.sendLinkGoogleAccountRedirect(w, r)
		}
	case VerificationGitHub:
		h.sendGithubActionRedirect(w, r, ActionLinkGoogleAccount, "/user/oauth/github/callback/action")
	default:
		w.WriteHeader(http.StatusBadRequest)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
	}
}

func (h *Handlers) handleUserSettingsLinkGithubAccountPost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	verificationMethod, ok := ParseVerificationMethod(data.Verification)
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	switch verificationMethod {
	case VerificationPassword:
		ok := h.validateAndCheckPasswordField(w, r, data, data.LinkGithub.Password, "LinkGithub.Password")
		if ok {
			h.sendLinkGithubAccountRedirect(w, r)
		}
	case VerificationGoogle:
		h.sendGoogleActionRedirect(w, r, ActionLinkGitHubAccount, "/user/oauth/google/callback/action")
	default:
		w.WriteHeader(http.StatusBadRequest)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
	}
}

func (h *Handlers) sendLinkGoogleAccountRedirect(w http.ResponseWriter, r *http.Request) {
	state := crypto.GenerateRandomString(30)
	c := &http.Cookie{
		Name:     "google_state",
		Value:    state,
		Path:     "/",
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	redirectUrl := h.GoogleOauthService.GetRedirectUrlWithCustomCallbackPath(state, "/user/oauth/google/callback/link")
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (h *Handlers) sendLinkGithubAccountRedirect(w http.ResponseWriter, r *http.Request) {
	state := crypto.GenerateRandomString(16)
	c := &http.Cookie{
		Name:     "github_state",
		Value:    state,
		Path:     "/",
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	redirectUrl := h.GithubOauthService.GetRedirectUrlWithCustomCallbackPath(state, "/user/oauth/github/callback/link")
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}
