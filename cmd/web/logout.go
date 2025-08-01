package main

import (
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserId")
	app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  app.getLocalizer(r).MustLocalize(&i18n.LocalizeConfig{MessageID: "toast.logout"}),
			Variant:  "simple",
			Duration: 2000,
		},
	})
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
