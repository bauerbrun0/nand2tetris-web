package main

import (
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/ui/pages/landingpage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/registerpage"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	msg := app.sessionManager.GetString(r.Context(), "message")
	app.logger.Info("getting message from session", "message", msg)
	app.sessionManager.Put(r.Context(), "message", "Hello World!")
	app.logger.Info("stored message to session")
	app.render(r.Context(), w, r, landingpage.Page())
}

func (app *application) userRegister(w http.ResponseWriter, r *http.Request) {
	app.render(r.Context(), w, r, registerpage.Page())
}
