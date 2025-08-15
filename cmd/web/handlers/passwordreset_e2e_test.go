package handlers_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserResetPasswordSendCode(t *testing.T) {
	ts, _, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	t.Run("Can visit page", func(t *testing.T) {
		result := ts.Get(t, "/user/reset-password/send-code")
		assert.Equal(t, http.StatusOK, result.Status)
		csrfToken := testutils.ExtractCSRFToken(t, result.Body)
		assert.NotEmpty(t, csrfToken)
	})

	t.Run("Redirect if already logged in", func(t *testing.T) {
		ts.MustLogIn(t, testutils.LoginParams{})
		result := ts.Get(t, "/user/reset-password/send-code")
		assert.Equal(t, http.StatusSeeOther, result.Status)
	})
}

func TestUserResetPasswordSendCodePost(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	result := ts.Get(t, "/user/reset-password/send-code")
	validCSRFToken := testutils.ExtractCSRFToken(t, result.Body)

	tests := []struct {
		name      string
		email     string
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
				testutils.ExpectGetUserByEmailReturnVerifiedEmailUser(t, queries)
				queries.EXPECT().InvalidatePasswordResetRequestsOfUser(t.Context(), mock.Anything).
					Return(nil).Once()
				queries.EXPECT().GetPasswordResetRequestByCode(t.Context(), mock.Anything).
					Return(models.PasswordResetRequest{}, pgx.ErrNoRows).Once()
				testutils.ExpectCreatePasswordResetRequestReturnsRequest(t, queries)
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
			name:      "Redirect if already logged in",
			email:     testutils.MockEmail,
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
			name:      "Empty email",
			email:     "",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "Invalid email",
			email:     "walter@",
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

			result := ts.PostForm(t, "/user/reset-password/send-code", form)
			assert.Equal(t, tt.wantCode, result.Status)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}

func TestUserResetPasswordEnterCode(t *testing.T) {
	ts, _, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	t.Run("Can visit page", func(t *testing.T) {
		result := ts.Get(t, "/user/reset-password/enter-code")
		assert.Equal(t, http.StatusOK, result.Status)
		csrfToken := testutils.ExtractCSRFToken(t, result.Body)
		assert.NotEmpty(t, csrfToken)
	})

	t.Run("Redirect if already logged in", func(t *testing.T) {
		ts.MustLogIn(t, testutils.LoginParams{})
		result := ts.Get(t, "/user/reset-password/enter-code")
		assert.Equal(t, http.StatusSeeOther, result.Status)
	})
}

func TestUserResetPasswordEnterCodePost(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	result := ts.Get(t, "/user/reset-password/enter-code")
	validCSRFToken := testutils.ExtractCSRFToken(t, result.Body)

	code := testutils.MockPasswordResetRequestCode

	tests := []struct {
		name      string
		email     string
		code      string
		csrfToken string
		wantCode  int
		before    func(t *testing.T)
		after     func(t *testing.T)
	}{
		{
			name:      "Valid submission",
			email:     testutils.MockEmail,
			code:      code,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				testutils.ExpectGetPasswordResetRequestByCodeReturnsRequest(t, queries)
			},
		},
		{
			name:      "Redirect if already logged in",
			email:     testutils.MockEmail,
			code:      code,
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
			name:      "Code expired",
			email:     testutils.MockEmail,
			code:      code,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
			before: func(t *testing.T) {
				testutils.ExpectGetPasswordResetRequestByCodeReturnsExpiredRequest(t, queries)
			},
		},
		{
			name:      "Code does not exists",
			email:     testutils.MockEmail,
			code:      code,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
			before: func(t *testing.T) {
				queries.EXPECT().GetPasswordResetRequestByCode(t.Context(), code).
					Return(models.PasswordResetRequest{}, pgx.ErrNoRows).Once()
			},
		},
		{
			name:      "Empty code",
			email:     testutils.MockEmail,
			code:      "",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "Empty email",
			email:     "",
			code:      code,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "Empty csrf token",
			email:     testutils.MockEmail,
			code:      code,
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
			form.Add("code", tt.code)

			result := ts.PostForm(t, "/user/reset-password/enter-code", form)
			assert.Equal(t, tt.wantCode, result.Status)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}

func TestUserResetPassword(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	t.Run("Redirect if code is not in sessions", func(t *testing.T) {
		result := ts.Get(t, "/user/reset-password")
		assert.Equal(t, http.StatusSeeOther, result.Status)
	})

	t.Run("Can visit page if code is in session", func(t *testing.T) {
		result := ts.Get(t, "/user/reset-password/enter-code")
		csrfToken := testutils.ExtractCSRFToken(t, result.Body)

		resetCode := "12345678"

		form := url.Values{}
		form.Add("csrf_token", csrfToken)
		form.Add("email", testutils.MockEmail)
		form.Add("code", resetCode)

		passwordResetRequestMockReturn := models.PasswordResetRequest{
			ID:     testutils.MockId,
			UserID: testutils.MockUserId,
			Email:  testutils.MockEmail,
			Code:   resetCode,
			VerifyEmailAfter: pgtype.Bool{
				Bool:  false,
				Valid: true,
			},
			Expiry: pgtype.Timestamptz{
				Time:  time.Now().Add(time.Hour),
				Valid: true,
			},
		}

		queries.EXPECT().GetPasswordResetRequestByCode(t.Context(), resetCode).
			Return(passwordResetRequestMockReturn, nil).Once()

		postResult := ts.PostForm(t, "/user/reset-password/enter-code", form)
		assert.Equal(t, http.StatusSeeOther, postResult.Status)

		result = ts.Get(t, "/user/reset-password")
		assert.Equal(t, http.StatusOK, result.Status)
	})

	t.Run("Redirect if already logged in", func(t *testing.T) {
		ts.MustLogIn(t, testutils.LoginParams{})
		result := ts.Get(t, "/user/reset-password")
		assert.Equal(t, http.StatusSeeOther, result.Status)
	})
}

func TestUserResetPasswordPost(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	result := ts.Get(t, "/user/reset-password/enter-code")
	csrfToken := testutils.ExtractCSRFToken(t, result.Body)

	resetCode := testutils.MockPasswordResetRequestCode

	form := url.Values{}
	form.Add("csrf_token", csrfToken)
	form.Add("email", testutils.MockEmail)
	form.Add("code", resetCode)

	tests := []struct {
		name                    string
		code                    string
		newPassword             string
		newPasswordConfirmation string
		csrfToken               string
		wantCode                int
		before                  func(t *testing.T)
		after                   func(t *testing.T)
	}{
		{
			name:                    "Valid submission",
			code:                    resetCode,
			newPassword:             testutils.MockPassword,
			newPasswordConfirmation: testutils.MockPassword,
			csrfToken:               csrfToken,
			wantCode:                http.StatusSeeOther,
			before: func(t *testing.T) {
				testutils.ExpectGetPasswordResetRequestByCodeReturnsRequest(t, queries)
				queries.EXPECT().InvalidatePasswordResetRequest(t.Context(), mock.Anything).
					Return(nil).Once()
				queries.EXPECT().ChangeUserPasswordHash(t.Context(), mock.Anything).
					Return(nil).Once()
			},
		},
		{
			name:                    "Expired code",
			code:                    resetCode,
			newPassword:             testutils.MockPassword,
			newPasswordConfirmation: testutils.MockPassword,
			csrfToken:               csrfToken,
			wantCode:                http.StatusUnauthorized,
			before: func(t *testing.T) {
				testutils.ExpectGetPasswordResetRequestByCodeReturnsExpiredRequest(t, queries)
			},
		},
		{
			name:                    "Redirect if already logged in",
			code:                    resetCode,
			newPassword:             testutils.MockPassword,
			newPasswordConfirmation: testutils.MockPassword,
			csrfToken:               csrfToken,
			wantCode:                http.StatusSeeOther,
			before: func(t *testing.T) {
				ts.MustLogIn(t, testutils.LoginParams{})
			},
			after: func(t *testing.T) {
				ts.RemoveCookie(t, "session")
			},
		},
		{
			name:                    "Empty code",
			code:                    "",
			newPassword:             testutils.MockPassword,
			newPasswordConfirmation: testutils.MockPassword,
			csrfToken:               csrfToken,
			wantCode:                http.StatusUnprocessableEntity,
		},
		{
			name:                    "Empty password",
			code:                    resetCode,
			newPassword:             "",
			newPasswordConfirmation: testutils.MockPassword,
			csrfToken:               csrfToken,
			wantCode:                http.StatusUnprocessableEntity,
		},
		{
			name:                    "Empty password confirmation",
			code:                    resetCode,
			newPassword:             testutils.MockPassword,
			newPasswordConfirmation: "",
			csrfToken:               csrfToken,
			wantCode:                http.StatusUnprocessableEntity,
		},
		{
			name:                    "Passwords do not match",
			code:                    resetCode,
			newPassword:             testutils.MockPassword,
			newPasswordConfirmation: testutils.MockPassword + "x",
			csrfToken:               csrfToken,
			wantCode:                http.StatusUnprocessableEntity,
		},
		{
			name:                    "Empty csrf token",
			code:                    resetCode,
			newPassword:             testutils.MockPassword,
			newPasswordConfirmation: testutils.MockPassword,
			csrfToken:               "",
			wantCode:                http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}

			form := url.Values{}
			form.Add("csrf_token", tt.csrfToken)
			form.Add("new-password", tt.newPassword)
			form.Add("new-password-confirmation", tt.newPasswordConfirmation)
			form.Add("code", tt.code)

			result := ts.PostForm(t, "/user/reset-password", form)
			assert.Equal(t, tt.wantCode, result.Status)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}
