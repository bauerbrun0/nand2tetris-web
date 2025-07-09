package main

import (
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/ui/pages/landingpage"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(r.Context(), w, r, landingpage.Page())
}
