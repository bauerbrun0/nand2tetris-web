package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/internal/crypto"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
)

func (h *Handlers) UserLoginGithub(w http.ResponseWriter, r *http.Request) {
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

	redirectUrl := h.GithubOauthService.GetRedirectUrl(state)
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (h *Handlers) UserLoginGithubCallback(w http.ResponseWriter, r *http.Request) {
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
		http.Redirect(w, r, "/user/login", http.StatusBadRequest)
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
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
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

	user, err := h.UserService.AuthenticateOAuthUser(oauthUser, models.ProviderGitHub)
	if err != nil && !errors.Is(err, services.ErrUserAlreadyExists) {
		h.ServerError(w, r, err)
		return
	}

	if err != nil && errors.Is(err, services.ErrUserAlreadyExists) {
		h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.oauth_user_exists"),
				Variant:  "error",
				Duration: 5000,
			},
		})
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	h.SessionManager.Put(r.Context(), "authenticatedUserId", user.ID)
	h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.TTemplate("toast.user_welcome", map[string]string{"Username": user.Username}),
			Variant:  "simple",
			Duration: 2000,
		},
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handlers) UserLoginGoogle(w http.ResponseWriter, r *http.Request) {
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

	redirectUrl := h.GoogleOauthService.GetRedirectUrl(state)
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (h *Handlers) UserLoginGoogleCallback(w http.ResponseWriter, r *http.Request) {
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
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
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
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := h.GoogleOauthService.ExchangeCodeForToken(services.TokenExchangeOptions{
		Code:         code,
		RedirectPath: "/user/oauth/google/callback/login",
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

	user, err := h.UserService.AuthenticateOAuthUser(oauthUser, models.ProviderGoogle)
	if err != nil && !errors.Is(err, services.ErrUserAlreadyExists) {
		h.ServerError(w, r, err)
		return
	}

	if err != nil && errors.Is(err, services.ErrUserAlreadyExists) {
		h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.oauth_user_exists"),
				Variant:  "error",
				Duration: 5000,
			},
		})
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	h.SessionManager.Put(r.Context(), "authenticatedUserId", user.ID)
	h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.TTemplate("toast.user_welcome", map[string]string{"Username": user.Username}),
			Variant:  "simple",
			Duration: 2000,
		},
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
