package main

import (
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/validator"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/usersettingspage"
)

func (app *application) userSettings(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	data := usersettingspage.UserSettingsPageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()
	app.render(r.Context(), w, r, usersettingspage.Page(data))
}

func (app *application) userSettingsPost(w http.ResponseWriter, r *http.Request) {
	basePageData := app.newPageData(r)
	data := usersettingspage.UserSettingsPageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()

	err := app.decodePostForm(r, &data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	switch data.Form {
	case "change-password":
		app.handleUserSettingsChangePasswordPost(w, r, &data)
	case "change-email":
		app.handleUserSettingsChangeEmailPost(w, r, &data)
	case "change-email-send-code":
		app.handleUserSettingsChangeEmailSendCodePost(w, r, &data)
	case "delete-account":
		app.handleUserSettingsDeleteAccountPost(w, r, &data)
	case "create-password":
		app.handleUserSettingsCreatePasswordPost(w, r, &data)
	case "link-google-account":
		app.handleUserSettingsLinkGoogleAccountPost(w, r, &data)
	case "link-github-account":
		app.handleUserSettingsLinkGithubAccountPost(w, r, &data)
	case "unlink-google-account":
		app.handleUserSettingsUnlinkGoogleAccountPost(w, r, &data)
	case "unlink-github-account":
		app.handleUserSettingsUnlinkGithubAccountPost(w, r, &data)
	default:
		w.WriteHeader(http.StatusBadRequest)
		app.render(r.Context(), w, r, usersettingspage.Page(data))
	}
}
