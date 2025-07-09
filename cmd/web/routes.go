package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static"))))

	mux.HandleFunc("GET /{$}", app.home)

	commonChain := alice.New(app.recoverPanic, app.logRequest)

	return commonChain.Then(mux)
}
