package testutils

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/application"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/middleware"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/routes"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/models/mocks"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/go-playground/form"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func NewTestApplication(
	t *testing.T, queries models.DBQueries, githubOauthService, googleOauthService services.OAuthService,
) *application.Application {
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	formDecoder := form.NewDecoder()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	bundle.LoadMessageFile("internal/translations/en.yaml")

	ctx := t.Context()
	txStarter := mocks.NewMockTxStarter(queries)

	emailSender := services.NewConsoleEmailSender(logger)
	emailService := services.NewEmailService(emailSender, logger)
	userService := services.NewUserService(logger, nil, queries, txStarter, ctx)

	cfg := application.Config{
		Env: "test",
	}

	return &application.Application{
		Logger:             logger,
		Config:             cfg,
		UserService:        userService,
		EmailService:       emailService,
		GithubOauthService: githubOauthService,
		GoogleOauthService: googleOauthService,
		SessionManager:     sessionManager,
		FormDecoder:        formDecoder,
		Bundle:             bundle,
	}
}

func NewTestHandlers(t *testing.T, app *application.Application) *handlers.Handlers {
	return handlers.NewHandlers(app)
}

func NewTestMiddleware(t *testing.T, app *application.Application) *middleware.Middleware {
	return middleware.NewMiddleware(app)
}

func NewTestRoutes(
	t *testing.T, app *application.Application, handlers *handlers.Handlers, middleware *middleware.Middleware,
) http.Handler {
	return routes.GetRoutes(app, middleware, handlers)
}

type testServer struct {
	*httptest.Server
}

func NewTestServer(t *testing.T, queries models.DBQueries, githubOauthService, googleOauthService services.OAuthService) *testServer {
	app := NewTestApplication(t, queries, githubOauthService, googleOauthService)
	h := NewTestHandlers(t, app)
	m := NewTestMiddleware(t, app)
	routes := NewTestRoutes(t, app, h, m)
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

func (ts *testServer) Get(t *testing.T, urlPath string) (int, http.Header, string) {
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
