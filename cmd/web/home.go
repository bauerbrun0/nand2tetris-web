package main

import (
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/ui/pages/homepage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/landingpage"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newPageData(r)
	if app.isAuthenticated(r) {
		app.render(r.Context(), w, r, homepage.Page(data))
		return
	}
	app.render(r.Context(), w, r, landingpage.Page(data))
}
