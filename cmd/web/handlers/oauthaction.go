package handlers

import (
	"net/http"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/internal/crypto"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
)

func (h *Handlers) sendGoogleActionRedirect(w http.ResponseWriter, r *http.Request, action Action, callbackPath string) {
	h.SessionManager.Put(r.Context(), "authenticated-action", action)

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

	redirectUrl := h.GoogleOauthService.GetRedirectUrlWithCustomCallbackPath(state, callbackPath)
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (h *Handlers) sendGithubActionRedirect(w http.ResponseWriter, r *http.Request, action Action, callbackPath string) {
	h.SessionManager.Put(r.Context(), "authenticated-action", action)

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

	redirectUrl := h.GithubOauthService.GetRedirectUrlWithCustomCallbackPath(state, callbackPath)
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (h *Handlers) UserGoogleActionCallback(w http.ResponseWriter, r *http.Request) {
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
		RedirectPath: "/user/oauth/google/callback/action",
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

	userId, ok, err := h.UserService.GetUserIdByUserProviderId(models.ProviderGoogle, oauthUser.Id)
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	if !ok {
		h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.oauth_user_not_exists"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
		return
	}

	if userId != data.UserInfo.ID {
		h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.oauth_user_not_current_user"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
		return
	}

	action, ok := h.SessionManager.Pop(r.Context(), "authenticated-action").(Action)
	if !ok {
		h.ServerError(w, r, ErrInvalidActionInSession)
		return
	}

	switch action {
	case ActionLinkGitHubAccount:
		h.sendLinkGithubAccountRedirect(w, r)
	case ActionUnlinkGoogleAccount:
		h.unlinkOAuthAccount(w, r, &data, models.ProviderGoogle)
	case ActionUnlinkGitHubAccount:
		h.unlinkOAuthAccount(w, r, &data, models.ProviderGitHub)
	case ActionDeleteAccount:
		h.deleteAccount(w, r)
	default:
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
	}
}

func (h *Handlers) UserGithubActionCallback(w http.ResponseWriter, r *http.Request) {
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

	userId, ok, err := h.UserService.GetUserIdByUserProviderId(models.ProviderGitHub, oauthUser.Id)
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	if !ok {
		h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.oauth_user_not_exists"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
		return
	}

	if userId != data.UserInfo.ID {
		h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.oauth_user_not_current_user"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
		return
	}

	action, ok := h.SessionManager.Pop(r.Context(), "authenticated-action").(Action)
	if !ok {
		h.ServerError(w, r, ErrInvalidActionInSession)
		return
	}

	switch action {
	case ActionLinkGoogleAccount:
		h.sendLinkGoogleAccountRedirect(w, r)
	case ActionUnlinkGoogleAccount:
		h.unlinkOAuthAccount(w, r, &data, models.ProviderGoogle)
	case ActionUnlinkGitHubAccount:
		h.unlinkOAuthAccount(w, r, &data, models.ProviderGitHub)
	case ActionDeleteAccount:
		h.deleteAccount(w, r)
	default:
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
	}
}
