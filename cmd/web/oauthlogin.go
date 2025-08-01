package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/internal/crypto"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
)

func (app *application) userLoginGithub(w http.ResponseWriter, r *http.Request) {
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

	redirectUrl := app.githubOauthService.GetRedirectUrl(state)
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (app *application) userLoginGithubCallback(w http.ResponseWriter, r *http.Request) {
	data := app.newPageData(r)

	state, err := r.Cookie("github_state")
	if err != nil {
		app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
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
		app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
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
	token, err := app.githubOauthService.ExchangeCodeForToken(services.TokenExchangeOptions{
		Code: code,
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	oauthUser, err := app.githubOauthService.GetUserInfo(token)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	user, err := app.userService.AuthenticateOAuthUser(oauthUser, models.ProviderGitHub)
	if err != nil && !errors.Is(err, services.ErrUserAlreadyExists) {
		app.serverError(w, r, err)
		return
	}

	if err != nil && errors.Is(err, services.ErrUserAlreadyExists) {
		app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.oauth_user_exists"),
				Variant:  "error",
				Duration: 5000,
			},
		})
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
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

func (app *application) userLoginGoogle(w http.ResponseWriter, r *http.Request) {
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

	redirectUrl := app.googleOauthService.GetRedirectUrl(state)
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (app *application) userLoginGoogleCallback(w http.ResponseWriter, r *http.Request) {
	data := app.newPageData(r)

	state, err := r.Cookie("google_state")
	if err != nil {
		app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
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
		app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
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
	token, err := app.googleOauthService.ExchangeCodeForToken(services.TokenExchangeOptions{
		Code:         code,
		RedirectPath: "/user/oauth/google/callback/login",
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	oauthUser, err := app.googleOauthService.GetUserInfo(token)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	user, err := app.userService.AuthenticateOAuthUser(oauthUser, models.ProviderGoogle)
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
