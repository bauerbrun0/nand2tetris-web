package handlers_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	t.Run("Can visit page", func(t *testing.T) {
		result := ts.Get(t, "/user/login")

		assert.Equal(t, http.StatusOK, result.Status, "status code should be 200 OK")

		csrfToken := testutils.ExtractCSRFToken(t, result.Body)
		assert.NotEmptyf(t, csrfToken, "csrfToken should not be empty")
	})

	t.Run("Redirect if already logged in", func(t *testing.T) {
		ts.MustLogIn(t, queries, testutils.LoginUser{
			Username: "walt",
			Email:    "walter.white@example.com",
			Password: "LosPollos321",
		})
		result := ts.Get(t, "/user/login")
		assert.Equal(t, http.StatusSeeOther, result.Status, "status code should be 303 See Other")
	})
}

func TestUserLoginPost(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	result := ts.Get(t, "/user/login")
	validCSRFToken := testutils.ExtractCSRFToken(t, result.Body)

	var (
		validUsername     = "walter"
		validEmail        = "walter.white@example.com"
		validPassword     = "LosPollos321"
		validPasswordHash = testutils.MustHashPassword(t, validPassword)
	)

	returnUser := models.User{
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
	}

	tests := []struct {
		name      string
		username  string
		password  string
		csrfToken string
		wantCode  int
		before    func(t *testing.T)
		after     func(t *testing.T)
	}{
		{
			name:      "Valid submission with email",
			username:  validEmail,
			password:  validPassword,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				queries.EXPECT().GetUserByUsernameOrEmail(t.Context(), validEmail).
					Return(returnUser, nil).Once()
			},
			after: func(t *testing.T) {
				ts.RemoveCookie(t, "session")
			},
		},
		{
			name:      "Valid submission with username",
			username:  validUsername,
			password:  validPassword,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				queries.EXPECT().GetUserByUsernameOrEmail(t.Context(), validUsername).
					Return(returnUser, nil).Once()
			},
			after: func(t *testing.T) {
				ts.RemoveCookie(t, "session")
			},
		},
		{
			name:      "Wrong password",
			username:  validUsername,
			password:  validPassword + "wrong",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnauthorized,
			before: func(t *testing.T) {
				queries.EXPECT().GetUserByUsernameOrEmail(t.Context(), validUsername).
					Return(returnUser, nil).Once()
			},
		},
		{
			name:      "Wrong username",
			username:  "wrong",
			password:  validPassword,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnauthorized,
			before: func(t *testing.T) {
				queries.EXPECT().GetUserByUsernameOrEmail(t.Context(), "wrong").
					Return(models.User{}, pgx.ErrNoRows).Once()
			},
		},
		{
			name:      "Redirect if already logged in",
			username:  validUsername,
			password:  validPassword,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				ts.MustLogIn(t, queries, testutils.LoginUser{
					Username: validUsername,
					Email:    validEmail,
					Password: validPassword,
				})
			},
			after: func(t *testing.T) {
				ts.RemoveCookie(t, "session")
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
			if tt.before != nil {
				tt.before(t)
			}

			form := url.Values{}
			form.Add("username", tt.username)
			form.Add("password", tt.password)
			form.Add("csrf_token", tt.csrfToken)

			result := ts.PostForm(t, "/user/login", form)
			assert.Equal(t, tt.wantCode, result.Status)
			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}
