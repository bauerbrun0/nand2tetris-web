package main

import (
	"net/http"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/internal/crypto"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
)

func (app *application) sendGoogleActionRedirect(w http.ResponseWriter, r *http.Request, action, callbackPath string) {
	app.sessionManager.Put(r.Context(), "authenticated-action", action)

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

	redirectUrl := app.googleOauthService.GetRedirectUrlWithCustomCallbackPath(state, callbackPath)
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (app *application) sendGithubActionRedirect(w http.ResponseWriter, r *http.Request, action, callbackPath string) {
	app.sessionManager.Put(r.Context(), "authenticated-action", action)

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

	redirectUrl := app.githubOauthService.GetRedirectUrlWithCustomCallbackPath(state, callbackPath)
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (app *application) userGoogleActionCallback(w http.ResponseWriter, r *http.Request) {
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
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
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
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := app.googleOauthService.ExchangeCodeForToken(services.TokenExchangeOptions{
		Code:         code,
		RedirectPath: "/user/oauth/google/callback/action",
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

	userId, ok, err := app.userService.GetUserIdByUserProviderId(models.ProviderGoogle, oauthUser.Id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !ok {
		app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
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
		app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.oauth_user_not_current_user"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
		return
	}

	action := app.sessionManager.PopString(r.Context(), "authenticated-action")

	switch action {
	case "link-github-account":
		app.sendLinkGithubAccountRedirect(w, r)
	case "unlink-google-account":
		app.unlinkOAuthAccount(w, r, &data, models.ProviderGoogle)
	case "unlink-github-account":
		app.unlinkOAuthAccount(w, r, &data, models.ProviderGitHub)
	case "delete-account":
		app.deleteAccount(w, r)
	default:
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
	}
}

func (app *application) userGithubActionCallback(w http.ResponseWriter, r *http.Request) {
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
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
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
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
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

	userId, ok, err := app.userService.GetUserIdByUserProviderId(models.ProviderGitHub, oauthUser.Id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !ok {
		app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
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
		app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.oauth_user_not_current_user"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
		return
	}

	action := app.sessionManager.PopString(r.Context(), "authenticated-action")

	switch action {
	case "link-google-account":
		app.sendLinkGoogleAccountRedirect(w, r)
	case "unlink-google-account":
		app.unlinkOAuthAccount(w, r, &data, models.ProviderGoogle)
	case "unlink-github-account":
		app.unlinkOAuthAccount(w, r, &data, models.ProviderGitHub)
	case "delete-account":
		app.deleteAccount(w, r)
	default:
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
	}
}
