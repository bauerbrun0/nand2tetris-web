package handlers_test

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/cmd/web/application"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/middleware"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/routes"
)

func newTestApplication(t *testing.T) *application.Application {
	return &application.Application{
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func newHandlers(t *testing.T, app *application.Application) *handlers.Handlers {
	return handlers.NewHandlers(app)
}

func newMiddleware(t *testing.T, app *application.Application) *middleware.Middleware {
	return middleware.NewMiddleware(app)
}

func newRoutes(t *testing.T, app *application.Application, handlers *handlers.Handlers, middleware *middleware.Middleware) http.Handler {
	return routes.GetRoutes(app, middleware, handlers)
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T) *testServer {
	app := newTestApplication(t)
	h := newHandlers(t, app)
	m := newMiddleware(t, app)
	routes := newRoutes(t, app, h, m)
	ts := httptest.NewTLSServer(routes)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	// disable redirect-following for the test server client
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}
