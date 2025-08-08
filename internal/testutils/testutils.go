package testutils

import (
	"bytes"
	"encoding/gob"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/application"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/middleware"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/routes"
	"github.com/bauerbrun0/nand2tetris-web/internal"
	"github.com/bauerbrun0/nand2tetris-web/internal/crypto"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/models/mocks"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/go-playground/form"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func NewTestApplication(
	t *testing.T, queries models.DBQueries, githubOauthService, googleOauthService services.OAuthService, logs bool,
) *application.Application {
	gob.Register([]pages.Toast{})
	gob.Register(handlers.Action(""))
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	formDecoder := form.NewDecoder()

	var logger *slog.Logger
	if logs {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	} else {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	_, err := bundle.LoadMessageFileFS(internal.TranslationFiles, "translations/en.yaml")
	if err != nil {
		t.Fatal(err)
	}

	ctx := t.Context()
	txStarter := mocks.NewMockTxStarter(queries)

	emailSender := services.NewConsoleEmailSender(logger)
	emailService := services.NewEmailService(emailSender, logger)
	userService := services.NewUserService(logger, emailService, queries, txStarter, ctx)

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

func NewTestServer(t *testing.T, queries models.DBQueries, githubOauthService, googleOauthService services.OAuthService, logs bool) *testServer {
	app := NewTestApplication(t, queries, githubOauthService, googleOauthService, logs)
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

func (ts *testServer) PostForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	req, err := http.NewRequest(http.MethodPost, ts.URL+urlPath, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Origin", ts.URL)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rs, err := ts.Client().Do(req)
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

// Removes the cookie with the provided name from the ookie jar
// if it exists.
func (ts *testServer) RemoveCookie(t *testing.T, name string) {
	serverUrl, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	cookies := ts.Client().Jar.Cookies(serverUrl)
	var newCookies []*http.Cookie
	for _, c := range cookies {
		if c.Name != name {
			newCookies = append(newCookies, c)
		}
	}

	// have to create a new cookie jar because ts.Client().Jar.SetCookies does nothing
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	jar.SetCookies(serverUrl, newCookies)

	ts.Client().Jar = jar
}

func (ts *testServer) GetCookies(t *testing.T) []*http.Cookie {
	serverUrl, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	return ts.Client().Jar.Cookies(serverUrl)
}

func (ts *testServer) MustLogIn(t *testing.T, queries *mocks.MockDBQueries, username string, email string, password string) {
	// remove session cookie if already logged in
	ts.RemoveCookie(t, "session")
	// visit login page and get csrf token
	code, _, body := ts.Get(t, "/user/login")
	assert.Equal(t, http.StatusOK, code, "status code should be 200 OK")
	csrfToken := ExtractCSRFToken(t, body)
	assert.NotEmptyf(t, csrfToken, "csrfToken should not be empty")

	var hasher crypto.PasswordHasher
	returnUser := models.User{
		ID:       1,
		Username: username,
		Email:    email,
		EmailVerified: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		PasswordHash: pgtype.Text{
			String: MustHashPassword(t, hasher, password),
			Valid:  true,
		},
		Created: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}
	queries.EXPECT().GetUserByUsernameOrEmail(t.Context(), username).
		Return(returnUser, nil).Once()

	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)
	form.Add("csrf_token", csrfToken)

	code, _, _ = ts.PostForm(t, "/user/login", form)
	assert.Equal(t, http.StatusSeeOther, code)

	queries.EXPECT().GetUserInfo(t.Context(), int32(1)).Return(models.UserInfo{
		ID:       1,
		Username: username,
		Email:    email,
		EmailVerified: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		Created: pgtype.Timestamptz{
			Time:  time.Now().Add(-time.Minute),
			Valid: true,
		},
		IsPasswordSet:  true,
		LinkedAccounts: []string{},
	}, nil)
}

func (ts *testServer) MustRegister(t *testing.T, queries *mocks.MockDBQueries, username, email, password, emailVerificationCode string) {
	// remove session cookie if already logged in
	ts.RemoveCookie(t, "session")
	// visit register page and get csrf token
	code, _, body := ts.Get(t, "/user/register")
	assert.Equal(t, http.StatusOK, code, "status code should be 200 OK")
	csrfToken := ExtractCSRFToken(t, body)
	assert.NotEmptyf(t, csrfToken, "csrfToken should not be empty")

	var hasher crypto.PasswordHasher
	returnUser := models.User{
		ID:       1,
		Username: username,
		Email:    email,
		EmailVerified: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		PasswordHash: pgtype.Text{
			String: MustHashPassword(t, hasher, password),
			Valid:  true,
		},
		Created: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	returnEmailVerificationRequest := models.EmailVerificationRequest{
		ID:     1,
		UserID: 1,
		Email:  email,
		Code:   emailVerificationCode,
		Expiry: pgtype.Timestamptz{
			Time:  time.Now().Add(time.Hour),
			Valid: true,
		},
	}

	queries.EXPECT().CreateNewUser(t.Context(), mock.Anything).
		Return(returnUser, nil).Once()
	queries.EXPECT().GetEmailVerificationRequestByCode(t.Context(), mock.Anything).
		Return(models.EmailVerificationRequest{}, pgx.ErrNoRows).Once()
	queries.EXPECT().CreateEmailVerificationRequest(t.Context(), mock.Anything).
		Return(returnEmailVerificationRequest, nil).Once()

	form := url.Values{}
	form.Add("username", username)
	form.Add("email", email)
	form.Add("password", password)
	form.Add("password-confirmation", password)
	form.Add("terms", "on")
	form.Add("csrf_token", csrfToken)

	code, _, _ = ts.PostForm(t, "/user/register", form)
	assert.Equal(t, http.StatusSeeOther, code)
}

var rxCSRF = regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.+?)">`)

func ExtractCSRFToken(t *testing.T, body string) string {
	matches := rxCSRF.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatalf("no csrf token found in body")
	}

	return matches[1]
}

func MustHashPassword(t *testing.T, hasher crypto.PasswordHasher, password string) string {
	t.Helper()
	hash, err := hasher.GenerateFromPassword(password, crypto.DefaultPasswordHashParams)
	if err != nil {
		t.Fatalf("error generating hash for password %q: %v", password, err)
	}
	return hash
}
