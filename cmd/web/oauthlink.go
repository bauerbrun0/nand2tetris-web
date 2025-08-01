package main

import (
	"net/http"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/internal/crypto"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/usersettingspage"
)

func (app *application) userLinkGithubCallback(w http.ResponseWriter, r *http.Request) {
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

	_, exists, err := app.userService.GetUserIdByUserProviderId(models.ProviderGitHub, oauthUser.Id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if exists {
		app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.oauth_user_already_linked"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
		return
	}

	err = app.userService.AddOAuthAuthorization(oauthUser.Id, data.UserInfo.ID, models.ProviderGitHub)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.TTemplate("toast.oauth_user_successfully_linked", map[string]string{"Provider": "GitHub"}),
			Variant:  "success",
			Duration: 3000,
		},
	})
	http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
}

func (app *application) userLinkGoogleCallback(w http.ResponseWriter, r *http.Request) {
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
		RedirectPath: "/user/oauth/google/callback/link",
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

	_, exists, err := app.userService.GetUserIdByUserProviderId(models.ProviderGoogle, oauthUser.Id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if exists {
		app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.oauth_user_already_linked"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
		return
	}

	err = app.userService.AddOAuthAuthorization(oauthUser.Id, data.UserInfo.ID, models.ProviderGoogle)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.TTemplate("toast.oauth_user_successfully_linked", map[string]string{"Provider": "Google"}),
			Variant:  "success",
			Duration: 3000,
		},
	})
	http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
}

func (app *application) handleUserSettingsLinkGoogleAccountPost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	switch data.Verification {
	case "":
		ok := app.validateAndCheckPasswordField(w, r, data, data.LinkGooglePassword, "link-google/password")
		if ok {
			app.sendLinkGoogleAccountRedirect(w, r)
		}
	case "github":
		app.sendGithubActionRedirect(w, r, "link-google-account", "/user/oauth/github/callback/action")
	default:
		w.WriteHeader(http.StatusBadRequest)
		app.render(r.Context(), w, r, usersettingspage.Page(*data))
	}
}

func (app *application) handleUserSettingsLinkGithubAccountPost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	switch data.Verification {
	case "":
		ok := app.validateAndCheckPasswordField(w, r, data, data.LinkGithubPassword, "link-github/password")
		if ok {
			app.sendLinkGithubAccountRedirect(w, r)
		}
	case "google":
		app.sendGoogleActionRedirect(w, r, "link-github-account", "/user/oauth/google/callback/action")
	default:
		w.WriteHeader(http.StatusBadRequest)
		app.render(r.Context(), w, r, usersettingspage.Page(*data))
	}
}

func (app *application) sendLinkGoogleAccountRedirect(w http.ResponseWriter, r *http.Request) {
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

	redirectUrl := app.googleOauthService.GetRedirectUrlWithCustomCallbackPath(state, "/user/oauth/google/callback/link")
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (app *application) sendLinkGithubAccountRedirect(w http.ResponseWriter, r *http.Request) {
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

	redirectUrl := app.githubOauthService.GetRedirectUrlWithCustomCallbackPath(state, "/user/oauth/github/callback/link")
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}
