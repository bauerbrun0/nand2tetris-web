package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static"))))

	mux.HandleFunc("GET /{$}", app.home)

	return mux
}
