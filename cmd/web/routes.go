package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static"))))

	dynamicChain := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamicChain.ThenFunc(app.home))

	commonChain := alice.New(app.recoverPanic, app.logRequest, app.disableCacheInDevMode, app.commonHeaders)

	return commonChain.Then(mux)
}
