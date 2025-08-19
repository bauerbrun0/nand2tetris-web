package userhandlers_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserVerifyEmail(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	t.Run("Can visit page", func(t *testing.T) {
		result := ts.Get(t, "/user/verify-email")

		assert.Equal(t, http.StatusOK, result.Status, "status code should be 200 OK")

		csrfToken := testutils.ExtractCSRFToken(t, result.Body)
		assert.NotEmptyf(t, csrfToken, "csrfToken should not be empty")
	})

	t.Run("Can visit page after registration", func(t *testing.T) {
		ts.MustRegister(
			t,
			queries,
			testutils.MockUsername,
			testutils.MockEmail,
			testutils.MockPassword,
			"12345678",
		)
		result := ts.Get(t, "/user/verify-email")

		assert.Equal(t, http.StatusOK, result.Status, "status code should be 200 OK")
		csrfToken := testutils.ExtractCSRFToken(t, result.Body)
		assert.NotEmptyf(t, csrfToken, "csrfToken should not be empty")
		assert.Containsf(t, result.Body, testutils.MockEmail, "body should contain the email address: %s", testutils.MockEmail)
	})

	t.Run("Redirect if already logged in", func(t *testing.T) {
		ts.MustLogIn(t, testutils.LoginParams{})
		result := ts.Get(t, "/user/verify-email")
		assert.Equal(t, http.StatusSeeOther, result.Status, "status code should be 303 See Other")
	})
}

func TestUserVerifyEmailPost(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	result := ts.Get(t, "/user/verify-email")
	validCSRFToken := testutils.ExtractCSRFToken(t, result.Body)

	emailVerificationCode := testutils.MockEmailVerificationRequestCode

	tests := []struct {
		name      string
		code      string
		csrfToken string
		wantCode  int
		before    func(t *testing.T)
		after     func(t *testing.T)
	}{
		{
			name:      "Valid submission",
			code:      emailVerificationCode,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				testutils.ExpectGetEmailVerificationRequestByCodeReturnsRequest(t, queries)
				queries.EXPECT().InvalidateEmailVerificationRequest(t.Context(), mock.Anything).
					Return(nil).Once()
				queries.EXPECT().VerifyUserEmail(t.Context(), testutils.MockId).
					Return(nil).Once()
			},
		},
		{
			name:      "Valid submission after registration",
			code:      emailVerificationCode,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				ts.MustRegister(
					t,
					queries,
					testutils.MockUsername,
					testutils.MockEmail,
					testutils.MockPassword,
					emailVerificationCode,
				)
				testutils.ExpectGetEmailVerificationRequestByCodeReturnsRequest(t, queries)
				queries.EXPECT().InvalidateEmailVerificationRequest(t.Context(), mock.Anything).
					Return(nil).Once()
				queries.EXPECT().VerifyUserEmail(t.Context(), testutils.MockUserId).
					Return(nil).Once()
			},
		},
		{
			name:      "Code expired",
			code:      emailVerificationCode,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
			before: func(t *testing.T) {
				testutils.ExpectGetEmailVerificationRequestByCodeReturnsExpiredRequest(t, queries)
			},
		},
		{
			name:      "Empty code",
			code:      "",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "Inavlid csrf token",
			code:      emailVerificationCode,
			csrfToken: "",
			wantCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}

			form := url.Values{}
			form.Add("csrf_token", tt.csrfToken)
			form.Add("code", tt.code)

			result := ts.PostForm(t, "/user/verify-email", form)
			assert.Equal(t, tt.wantCode, result.Status)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}

func TestUserVerifyEmailResendCode(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	emailVerificationCode := testutils.MockEmailVerificationRequestCode

	t.Run("Can visit page", func(t *testing.T) {
		result := ts.Get(t, "/user/verify-email/send-code")

		assert.Equal(t, http.StatusOK, result.Status, "status code should be 200 OK")

		csrfToken := testutils.ExtractCSRFToken(t, result.Body)
		assert.NotEmptyf(t, csrfToken, "csrfToken should not be empty")
	})

	t.Run("Can visit page after registration", func(t *testing.T) {
		ts.MustRegister(
			t,
			queries,
			testutils.MockUsername,
			testutils.MockEmail,
			testutils.MockPassword,
			emailVerificationCode,
		)
		result := ts.Get(t, "/user/verify-email/send-code")

		assert.Equal(t, http.StatusOK, result.Status, "status code should be 200 OK")
		csrfToken := testutils.ExtractCSRFToken(t, result.Body)
		assert.NotEmptyf(t, csrfToken, "csrfToken should not be empty")
	})

	t.Run("Redirect if already logged in", func(t *testing.T) {
		ts.MustLogIn(t, testutils.LoginParams{})
		result := ts.Get(t, "/user/verify-email/send-code")
		assert.Equal(t, http.StatusSeeOther, result.Status, "status code should be 303 See Other")
	})
}

func TestUserVerifyEmailResendCodePost(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	result := ts.Get(t, "/user/verify-email/send-code")
	validCSRFToken := testutils.ExtractCSRFToken(t, result.Body)

	tests := []struct {
		email     string
		name      string
		csrfToken string
		wantCode  int
		before    func(t *testing.T)
		after     func(t *testing.T)
	}{
		{
			name:      "Valid submission",
			email:     testutils.MockEmail,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				testutils.ExpectGetUserByEmailReturnUnverifiedEmailUser(t, queries)
				queries.EXPECT().InvalidateEmailVerificationRequestsOfUser(t.Context(), mock.Anything).
					Return(nil).Once()
				queries.EXPECT().GetEmailVerificationRequestByCode(t.Context(), mock.Anything).
					Return(models.EmailVerificationRequest{}, pgx.ErrNoRows).Once()
				testutils.ExpectCreateEmailVerificationRequestReturnsRequest(t, queries)
			},
		},
		{
			name:      "No user with email",
			email:     testutils.MockEmail,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				queries.EXPECT().GetUserByEmail(t.Context(), testutils.MockEmail).
					Return(models.User{}, pgx.ErrNoRows).Once()
			},
		},
		{
			name:      "Email already verified",
			email:     testutils.MockEmail,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				testutils.ExpectGetUserByEmailReturnVerifiedEmailUser(t, queries)
			},
		},
		{
			name:      "Empty email",
			email:     "",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "Invalid email",
			email:     "walt@",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "Empty csrf token",
			email:     testutils.MockEmail,
			csrfToken: "",
			wantCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}

			form := url.Values{}
			form.Add("csrf_token", tt.csrfToken)
			form.Add("email", tt.email)

			result := ts.PostForm(t, "/user/verify-email/send-code", form)
			assert.Equal(t, tt.wantCode, result.Status)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}
