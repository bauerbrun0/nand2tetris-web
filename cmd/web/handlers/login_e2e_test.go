package handlers_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/internal/crypto"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	modelsmocks "github.com/bauerbrun0/nand2tetris-web/internal/models/mocks"
	servicemocks "github.com/bauerbrun0/nand2tetris-web/internal/services/mocks"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {
	githubOauthService := servicemocks.NewMockOAuthService(t)
	googleOauthService := servicemocks.NewMockOAuthService(t)
	queries := modelsmocks.NewMockDBQueries(t)

	ts := testutils.NewTestServer(t, queries, githubOauthService, googleOauthService, false)
	defer ts.Close()

	code, _, body := ts.Get(t, "/user/login")

	assert.Equal(t, http.StatusOK, code, "status code should be 200 OK")

	csrfToken := testutils.ExtractCSRFToken(t, body)
	assert.NotEmptyf(t, csrfToken, "csrfToken should not be empty")
}

func TestUserLoginPost(t *testing.T) {
	githubOauthService := servicemocks.NewMockOAuthService(t)
	googleOauthService := servicemocks.NewMockOAuthService(t)
	queries := modelsmocks.NewMockDBQueries(t)

	ts := testutils.NewTestServer(t, queries, githubOauthService, googleOauthService, false)
	defer ts.Close()

	_, _, body := ts.Get(t, "/user/login")
	validCSRFToken := testutils.ExtractCSRFToken(t, body)

	var hasher crypto.PasswordHasher

	var (
		validUsername     = "walter"
		validEmail        = "walter.white@example.com"
		validPassword     = "LosPollos321"
		validPasswordHash = testutils.MustHashPassword(t, hasher, validPassword)
	)

	tests := []struct {
		name        string
		username    string
		password    string
		csrfToken   string
		wantCode    int
		prepareMock func()
	}{
		{
			name:      "Valid submission with email",
			username:  validEmail,
			password:  validPassword,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			prepareMock: func() {
				queries.EXPECT().GetUserByUsernameOrEmail(t.Context(), validEmail).Return(models.User{
					ID:       1,
					Username: validUsername,
					Email:    validEmail,
					EmailVerified: pgtype.Bool{
						Bool:  true,
						Valid: true,
					},
					PasswordHash: pgtype.Text{
						String: validPasswordHash,
						Valid:  true,
					},
					Created: pgtype.Timestamptz{
						Time:  time.Now(),
						Valid: true,
					},
				}, nil).Once()
			},
		},
		{
			name:      "Valid submission with username",
			username:  validUsername,
			password:  validPassword,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			prepareMock: func() {
				queries.EXPECT().GetUserByUsernameOrEmail(t.Context(), validUsername).Return(models.User{
					ID:       1,
					Username: validUsername,
					Email:    validEmail,
					EmailVerified: pgtype.Bool{
						Bool:  true,
						Valid: true,
					},
					PasswordHash: pgtype.Text{
						String: validPasswordHash,
						Valid:  true,
					},
					Created: pgtype.Timestamptz{
						Time:  time.Now(),
						Valid: true,
					},
				}, nil).Once()
			},
		},
		{
			name:      "Wrong password",
			username:  validUsername,
			password:  validPassword + "wrong",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnauthorized,
			prepareMock: func() {
				queries.EXPECT().GetUserByUsernameOrEmail(t.Context(), validUsername).Return(models.User{
					ID:       1,
					Username: validUsername,
					Email:    validEmail,
					EmailVerified: pgtype.Bool{
						Bool:  true,
						Valid: true,
					},
					PasswordHash: pgtype.Text{
						String: validPasswordHash,
						Valid:  true,
					},
					Created: pgtype.Timestamptz{
						Time:  time.Now(),
						Valid: true,
					},
				}, nil).Once()
			},
		},
		{
			name:      "Wrong username",
			username:  "wrong",
			password:  validPassword,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnauthorized,
			prepareMock: func() {
				queries.EXPECT().GetUserByUsernameOrEmail(t.Context(), "wrong").
					Return(models.User{}, pgx.ErrNoRows).Once()
			},
		},
		{
			name:      "Invalid csrf token",
			username:  validUsername,
			password:  validPassword,
			csrfToken: validCSRFToken + "wrong",
			wantCode:  http.StatusBadRequest,
		},
		{
			name:      "empty username",
			username:  "",
			password:  validPassword,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "empty password",
			username:  validUsername,
			password:  "",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "too long password",
			username:  validUsername,
			password:  strings.Repeat("x", 100),
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepareMock != nil {
				tt.prepareMock()
			}
			ts.RemoveCookie(t, "session")

			form := url.Values{}
			form.Add("username", tt.username)
			form.Add("password", tt.password)
			form.Add("csrf_token", tt.csrfToken)

			code, _, _ := ts.PostForm(t, "/user/login", form)
			assert.Equal(t, tt.wantCode, code)
		})
	}

}
