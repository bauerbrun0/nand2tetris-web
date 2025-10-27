package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"maps"
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
		ShowFooter:      true,
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

func (app *Application) WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	jsonData = append(jsonData, '\n') // looks better in terminal

	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonData)
	return nil
}

func (app *Application) WriteJSONError(w http.ResponseWriter, r *http.Request, status int, message any) {
	err := app.WriteJSON(w, status, map[string]any{"error": message}, nil)
	if err != nil {
		app.Logger.Error(
			"Failed to write JSON error",
			slog.String("json_write_error", err.Error()),
			slog.Any("error_message", message),
			slog.String("method", r.Method),
			slog.String("uri", r.URL.RequestURI()),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *Application) WriteJSONServerError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.Logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri))
	app.WriteJSONError(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (app *Application) WriteJSONNotFoundError(w http.ResponseWriter, r *http.Request) {
	app.WriteJSONError(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func (app *Application) WriteJSONMethodNotAllowedError(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not allowed for this resource", r.Method)
	app.WriteJSONError(w, r, http.StatusMethodNotAllowed, message)
}

func (app *Application) WriteJSONBadRequestError(w http.ResponseWriter, r *http.Request, error string) {
	app.WriteJSONError(w, r, http.StatusBadRequest, error)
}

func (app *Application) ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	return nil
}
