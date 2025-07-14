package main

import (
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	var fs http.Handler

	if app.config.env == "production" {
		fs = http.FileServerFS(ui.StaticFiles)
	} else {
		fs = http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static")))
	}

	mux.Handle("GET /static/", app.disableCacheInDevMode(fs))

	dynamicChain := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamicChain.ThenFunc(app.home))
	mux.Handle("GET /user/register", dynamicChain.ThenFunc(app.userRegister))
	mux.Handle("POST /user/register", dynamicChain.ThenFunc(app.userRegisterPost))

	commonChain := alice.New(app.recoverPanic, app.logRequest, app.commonHeaders)

	return commonChain.Then(mux)
}
