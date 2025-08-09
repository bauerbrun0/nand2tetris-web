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
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/application"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/middleware"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/routes"
	"github.com/bauerbrun0/nand2tetris-web/internal"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/models/mocks"
	modelsmocks "github.com/bauerbrun0/nand2tetris-web/internal/models/mocks"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	servicemocks "github.com/bauerbrun0/nand2tetris-web/internal/services/mocks"
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
	txStarter := modelsmocks.NewMockTxStarter(queries)

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

type TestServerOptions struct {
	Logs bool
}

func NewTestServer(t *testing.T, options TestServerOptions) (
	ts *testServer, queries *mocks.MockDBQueries, githubOauthService, googleOauthService *servicemocks.MockOAuthService,
) {
	githubOauthService = servicemocks.NewMockOAuthService(t)
	googleOauthService = servicemocks.NewMockOAuthService(t)
	queries = modelsmocks.NewMockDBQueries(t)

	app := NewTestApplication(t, queries, githubOauthService, googleOauthService, options.Logs)
	h := NewTestHandlers(t, app)
	m := NewTestMiddleware(t, app)
	routes := NewTestRoutes(t, app, h, m)
	server := httptest.NewTLSServer(routes)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	server.Client().Jar = jar

	// disable redirect-following for the test server client
	server.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{Server: server}, queries, githubOauthService, googleOauthService
}

type GetResult struct {
	Status      int
	Header      http.Header
	Body        string
	RedirectUrl *url.URL
}

func (ts *testServer) Get(t *testing.T, urlPath string) GetResult {
	var redirectUrl *url.URL
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		redirectUrl = req.URL
		return http.ErrUseLastResponse
	}
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	rsBody, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(rsBody)

	return GetResult{
		Status:      rs.StatusCode,
		Header:      rs.Header,
		Body:        string(rsBody),
		RedirectUrl: redirectUrl,
	}
}

type PostFormResult struct {
	Status      int
	Header      http.Header
	Body        string
	RedirectURL *url.URL
}

func (ts *testServer) PostForm(t *testing.T, urlPath string, form url.Values) PostFormResult {
	req, err := http.NewRequest(http.MethodPost, ts.URL+urlPath, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Origin", ts.URL)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var redirectUrl *url.URL
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		redirectUrl = req.URL
		return http.ErrUseLastResponse
	}

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

	return PostFormResult{
		Status:      rs.StatusCode,
		Header:      rs.Header,
		Body:        string(body),
		RedirectURL: redirectUrl,
	}
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

type LoginUser struct {
	Username       string
	Email          string
	Password       string
	LinkedAccounts []models.Provider
}

func (ts *testServer) MustLogIn(t *testing.T, queries *mocks.MockDBQueries, user LoginUser) (passwordHash string) {
	// remove session cookie if already logged in
	ts.RemoveCookie(t, "session")

	// visit login page and get csrf token
	result := ts.Get(t, "/user/login")
	assert.Equal(t, http.StatusOK, result.Status, "status code should be 200 OK")
	csrfToken := ExtractCSRFToken(t, result.Body)
	assert.NotEmptyf(t, csrfToken, "csrfToken should not be empty")

	if user.Password != "" {
		passwordHash = MustHashPassword(t, user.Password)
	} else {
		passwordHash = ""
	}

	returnUser := models.User{
		ID:       1,
		Username: user.Username,
		Email:    user.Email,
		EmailVerified: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		PasswordHash: pgtype.Text{
			String: passwordHash,
			Valid:  true,
		},
		Created: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}
	queries.EXPECT().GetUserByUsernameOrEmail(t.Context(), user.Username).
		Return(returnUser, nil).Once()

	form := url.Values{}
	form.Add("username", user.Username)
	form.Add("password", user.Password)
	form.Add("csrf_token", csrfToken)

	postResult := ts.PostForm(t, "/user/login", form)
	assert.Equal(t, http.StatusSeeOther, postResult.Status)

	stringLinkedAccounts := []string{}
	for _, a := range user.LinkedAccounts {
		stringLinkedAccounts = append(stringLinkedAccounts, string(a))
	}

	queries.EXPECT().GetUserInfo(t.Context(), int32(1)).Return(models.UserInfo{
		ID:       1,
		Username: user.Username,
		Email:    user.Email,
		EmailVerified: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		Created: pgtype.Timestamptz{
			Time:  time.Now().Add(-time.Minute),
			Valid: true,
		},
		IsPasswordSet:  passwordHash != "",
		LinkedAccounts: stringLinkedAccounts,
	}, nil)

	return passwordHash
}

type UserSettingsOAuthActionParams struct {
	Action       handlers.Action
	Verification handlers.VerificationMethod
	CSRFToken    string
	FormData     map[string]string
}

func (ts *testServer) MustSendUserSettingsOAuthAction(t *testing.T, githubOauthService, googleOauthService *servicemocks.MockOAuthService, params UserSettingsOAuthActionParams) (generatedState string) {
	switch params.Verification {
	case handlers.VerificationGitHub:
		githubOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, mock.Anything).RunAndReturn(func(state string, callbackPath string) string {
			generatedState = state
			return "https://github.com/login/oauth/authorize"
		}).Once()
	case handlers.VerificationGoogle:
		googleOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, mock.Anything).RunAndReturn(func(state string, callbackPath string) string {
			generatedState = state
			return "https://accounts.google.com/o/oauth2/v2/auth"
		}).Once()
	default:
		t.Fatalf("Unexpected verification method: %v", params.Verification)
	}

	form := url.Values{}
	for key, value := range params.FormData {
		form.Add(key, value)
	}
	form.Add("csrf_token", params.CSRFToken)
	form.Add("Action", string(params.Action))
	form.Add("Verification", string(params.Verification))

	result := ts.PostForm(t, "/user/settings", form)
	assert.Equal(t, http.StatusSeeOther, result.Status)
	assert.NotEmpty(t, generatedState)
	return generatedState
}

func (ts *testServer) MustRegister(t *testing.T, queries *mocks.MockDBQueries, username, email, password, emailVerificationCode string) {
	// remove session cookie if already logged in
	ts.RemoveCookie(t, "session")
	// visit register page and get csrf token
	result := ts.Get(t, "/user/register")
	assert.Equal(t, http.StatusOK, result.Status, "status code should be 200 OK")
	csrfToken := ExtractCSRFToken(t, result.Body)
	assert.NotEmptyf(t, csrfToken, "csrfToken should not be empty")

	returnUser := models.User{
		ID:       1,
		Username: username,
		Email:    email,
		EmailVerified: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		PasswordHash: pgtype.Text{
			String: MustHashPassword(t, password),
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

	postResult := ts.PostForm(t, "/user/register", form)
	assert.Equal(t, http.StatusSeeOther, postResult.Status)
}
