package application

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/bauerbrun0/nand2tetris-web/internal/appctx"
	"github.com/bauerbrun0/nand2tetris-web/internal/ctxi18n"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/go-playground/form"
	"github.com/justinas/nosurf"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func (app *Application) Render(ctx context.Context, w http.ResponseWriter, r *http.Request, t templ.Component) {
	err := t.Render(ctx, w)
	if err != nil {
		app.ServerError(w, r, err)
	}
}

func (app *Application) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.Logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) DecodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.FormDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return nil
}

func (app *Application) NewPageData(r *http.Request) pages.PageData {
	return pages.PageData{
		CurrentYear:     time.Now().Year(),
		IsAuthenticated: app.IsAuthenticated(r),
		UserInfo:        app.GetAuthenticatedUserInfo(r),
		InitialToasts:   app.GetInitialToasts(r),
		Localizer:       app.GetLocalizer(r),
		CSRFToken:       nosurf.Token(r),
		SveltePage:      pages.SveltePageNone,
	}
}

func (app *Application) IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(appctx.IsAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

func (app *Application) GetAuthenticatedUserInfo(r *http.Request) *pages.UserInfo {
	user, ok := r.Context().Value(appctx.AuthenticatedUserInfoKey).(*pages.UserInfo)
	if !ok {
		return nil
	}
	return user
}

func (app *Application) GetInitialToasts(r *http.Request) []pages.Toast {
	toasts, ok := r.Context().Value(appctx.InitialToastsKey).([]pages.Toast)
	if !ok {
		return []pages.Toast{}
	}
	return toasts
}

func (app *Application) GetLocalizer(r *http.Request) *i18n.Localizer {
	return ctxi18n.GetLocalizer(r.Context())
}
