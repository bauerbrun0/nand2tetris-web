package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/bauerbrun0/nand2tetris-web/internal/appctx"
	"github.com/bauerbrun0/nand2tetris-web/internal/ctxi18n"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/go-playground/form"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return nil
}

func (app *application) newPageData(r *http.Request) pages.PageData {
	return pages.PageData{
		CurrentYear:     time.Now().Year(),
		IsAuthenticated: app.isAuthenticated(r),
		UserInfo:        app.getAuthenticatedUserInfo(r),
		InitialToasts:   app.getInitialToasts(r),
		Localizer:       app.getLocalizer(r),
	}
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(appctx.IsAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

func (app *application) getAuthenticatedUserInfo(r *http.Request) *models.GetUserInfoRow {
	user, ok := r.Context().Value(appctx.AuthenticatedUserInfoKey).(*models.GetUserInfoRow)
	if !ok {
		return nil
	}
	return user
}

func (app *application) getInitialToasts(r *http.Request) []pages.Toast {
	toasts, ok := r.Context().Value(appctx.InitialToastsKey).([]pages.Toast)
	if !ok {
		return []pages.Toast{}
	}
	return toasts
}

func (app *application) getLocalizer(r *http.Request) *i18n.Localizer {
	return ctxi18n.GetLocalizer(r.Context())
}
