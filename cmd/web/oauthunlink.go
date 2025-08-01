package main

import (
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/usersettingspage"
)

func (app *application) unlinkOAuthAccount(w http.ResponseWriter, r *http.Request, data *pages.PageData, provider models.Provider) {
	if len(data.UserInfo.LinkedAccounts) == 1 && !data.UserInfo.IsPasswordSet {
		app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
			{
				Message:  data.T("toast.only_login_method"),
				Variant:  "error",
				Duration: 3000,
			},
		})
		http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
		return
	}

	err := app.userService.RemoveOAuthAuthorization(data.UserInfo.ID, provider)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.TTemplate("toast.oauth_user_successfully_unlinked", map[string]string{"Provider": string(provider)}),
			Variant:  "success",
			Duration: 3000,
		},
	})
	http.Redirect(w, r, "/user/settings", http.StatusSeeOther)
}

func (app *application) handleUserSettingsUnlinkGoogleAccountPost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	switch data.Verification {
	case "":
		ok := app.validateAndCheckPasswordField(w, r, data, data.UnlinkGooglePassword, "unlink-google/password")
		if ok {
			newPageData := app.newPageData(r)
			app.unlinkOAuthAccount(w, r, &newPageData, models.ProviderGoogle)
		}
	case "google":
		app.sendGoogleActionRedirect(w, r, "unlink-google-account", "/user/oauth/google/callback/action")
	case "github":
		app.sendGithubActionRedirect(w, r, "unlink-google-account", "/user/oauth/github/callback/action")
	default:
		w.WriteHeader(http.StatusBadRequest)
		app.render(r.Context(), w, r, usersettingspage.Page(*data))
	}
}

func (app *application) handleUserSettingsUnlinkGithubAccountPost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	switch data.Verification {
	case "":
		ok := app.validateAndCheckPasswordField(w, r, data, data.UnlinkGithubPassword, "unlink-github/password")
		if ok {
			newPageData := app.newPageData(r)
			app.unlinkOAuthAccount(w, r, &newPageData, models.ProviderGitHub)
		}
	case "google":
		app.sendGoogleActionRedirect(w, r, "unlink-github-account", "/user/oauth/google/callback/action")
	case "github":
		app.sendGithubActionRedirect(w, r, "unlink-github-account", "/user/oauth/github/callback/action")
	default:
		w.WriteHeader(http.StatusBadRequest)
		app.render(r.Context(), w, r, usersettingspage.Page(*data))
	}
}
