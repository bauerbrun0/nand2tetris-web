package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

func (app *application) render(ctx context.Context, w http.ResponseWriter, r *http.Request, t templ.Component) {
	err := t.Render(ctx, w)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
