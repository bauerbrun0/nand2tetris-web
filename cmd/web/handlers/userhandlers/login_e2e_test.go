package userhandlers_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {
	ts, _, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
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
		ts.MustLogIn(t, testutils.LoginParams{})
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
			username:  testutils.MockEmail,
			password:  testutils.MockPassword,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				testutils.ExpectGetUserByUsernameOrEmailReturnsUser(t, queries, testutils.MockEmail)
			},
			after: func(t *testing.T) {
				ts.RemoveCookie(t, "session")
			},
		},
		{
			name:      "Valid submission with username",
			username:  testutils.MockUsername,
			password:  testutils.MockPassword,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				testutils.ExpectGetUserByUsernameOrEmailReturnsUser(t, queries, testutils.MockUsername)
			},
			after: func(t *testing.T) {
				ts.RemoveCookie(t, "session")
			},
		},
		{
			name:      "Wrong password",
			username:  testutils.MockUsername,
			password:  testutils.MockPassword + "wrong",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnauthorized,
			before: func(t *testing.T) {
				testutils.ExpectGetUserByUsernameOrEmailReturnsUser(t, queries, testutils.MockUsername)
			},
		},
		{
			name:      "Wrong username",
			username:  "wrong",
			password:  testutils.MockPassword,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnauthorized,
			before: func(t *testing.T) {
				queries.EXPECT().GetUserByUsernameOrEmail(t.Context(), "wrong").
					Return(models.User{}, pgx.ErrNoRows).Once()
			},
		},
		{
			name:      "Redirect if already logged in",
			username:  testutils.MockUsername,
			password:  testutils.MockPassword,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				ts.MustLogIn(t, testutils.LoginParams{})
			},
			after: func(t *testing.T) {
				ts.RemoveCookie(t, "session")
			},
		},
		{
			name:      "Invalid csrf token",
			username:  testutils.MockUsername,
			password:  testutils.MockPassword,
			csrfToken: validCSRFToken + "wrong",
			wantCode:  http.StatusBadRequest,
		},
		{
			name:      "empty username",
			username:  "",
			password:  testutils.MockPassword,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "empty password",
			username:  testutils.MockUsername,
			password:  "",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "too long password",
			username:  testutils.MockUsername,
			password:  testutils.MockLongPassword,
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
