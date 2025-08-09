package handlers_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/internal/crypto"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserVerifyEmail(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	var (
		username              = "walter"
		email                 = "walter.white@example.com"
		password              = "LosPollos321"
		emailVerificationCode = "12345678"
	)

	t.Run("Can visit page", func(t *testing.T) {
		result := ts.Get(t, "/user/verify-email")

		assert.Equal(t, http.StatusOK, result.Status, "status code should be 200 OK")

		csrfToken := testutils.ExtractCSRFToken(t, result.Body)
		assert.NotEmptyf(t, csrfToken, "csrfToken should not be empty")
	})

	t.Run("Can visit page after registration", func(t *testing.T) {
		ts.MustRegister(t, queries, username, email, password, emailVerificationCode)
		result := ts.Get(t, "/user/verify-email")

		assert.Equal(t, http.StatusOK, result.Status, "status code should be 200 OK")
		csrfToken := testutils.ExtractCSRFToken(t, result.Body)
		assert.NotEmptyf(t, csrfToken, "csrfToken should not be empty")
		assert.Containsf(t, result.Body, email, "body should contain the email address: %s", email)
	})

	t.Run("Redirect if already logged in", func(t *testing.T) {
		ts.MustLogIn(t, queries, testutils.LoginUser{
			Username: username,
			Email:    email,
			Password: password,
		})
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

	var (
		username              = "walter"
		email                 = "walter.white@example.com"
		password              = "LosPollos321"
		emailVerificationCode = "12345678"
	)

	validEmailVerificationRequest := models.EmailVerificationRequest{
		ID:     1,
		UserID: 1,
		Email:  email,
		Code:   emailVerificationCode,
		Expiry: pgtype.Timestamptz{
			Time:  time.Now().Add(time.Hour),
			Valid: true,
		},
	}

	expiredEmailVerificationRequest := validEmailVerificationRequest
	expiredEmailVerificationRequest.Expiry.Time = time.Now().Add(-time.Hour)

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
				queries.EXPECT().GetEmailVerificationRequestByCode(t.Context(), emailVerificationCode).
					Return(validEmailVerificationRequest, nil).Once()
				queries.EXPECT().InvalidateEmailVerificationRequest(t.Context(), mock.Anything).
					Return(nil).Once()
				queries.EXPECT().VerifyUserEmail(t.Context(), int32(1)).
					Return(nil).Once()
			},
		},
		{
			name:      "Valid submission after registration",
			code:      emailVerificationCode,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				ts.MustRegister(t, queries, username, email, password, emailVerificationCode)
				queries.EXPECT().GetEmailVerificationRequestByCode(t.Context(), emailVerificationCode).
					Return(validEmailVerificationRequest, nil).Once()
				queries.EXPECT().InvalidateEmailVerificationRequest(t.Context(), mock.Anything).
					Return(nil).Once()
				queries.EXPECT().VerifyUserEmail(t.Context(), int32(1)).
					Return(nil).Once()
			},
		},
		{
			name:      "Code expired",
			code:      emailVerificationCode,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
			before: func(t *testing.T) {
				queries.EXPECT().GetEmailVerificationRequestByCode(t.Context(), emailVerificationCode).
					Return(expiredEmailVerificationRequest, nil).Once()
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

	var (
		username              = "walter"
		email                 = "walter.white@example.com"
		password              = "LosPollos321"
		emailVerificationCode = "12345678"
	)

	t.Run("Can visit page", func(t *testing.T) {
		result := ts.Get(t, "/user/verify-email/send-code")

		assert.Equal(t, http.StatusOK, result.Status, "status code should be 200 OK")

		csrfToken := testutils.ExtractCSRFToken(t, result.Body)
		assert.NotEmptyf(t, csrfToken, "csrfToken should not be empty")
	})

	t.Run("Can visit page after registration", func(t *testing.T) {
		ts.MustRegister(t, queries, username, email, password, emailVerificationCode)
		result := ts.Get(t, "/user/verify-email/send-code")

		assert.Equal(t, http.StatusOK, result.Status, "status code should be 200 OK")
		csrfToken := testutils.ExtractCSRFToken(t, result.Body)
		assert.NotEmptyf(t, csrfToken, "csrfToken should not be empty")
	})

	t.Run("Redirect if already logged in", func(t *testing.T) {
		ts.MustLogIn(t, queries, testutils.LoginUser{
			Username: username,
			Email:    email,
			Password: password,
		})
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

	var (
		email    = "walter.white@example.com"
		username = "walt"
		password = "LosPollos321"
	)

	var hasher crypto.PasswordHasher
	validUserMockReturn := models.User{
		ID:       1,
		Username: username,
		Email:    email,
		EmailVerified: pgtype.Bool{
			Bool:  false,
			Valid: true,
		},
		PasswordHash: pgtype.Text{
			String: testutils.MustHashPassword(t, hasher, password),
			Valid:  true,
		},
		Created: pgtype.Timestamptz{
			Time:  time.Now().Add(-time.Minute),
			Valid: true,
		},
	}

	emailVerifiedUserMockReturn := validUserMockReturn
	emailVerifiedUserMockReturn.EmailVerified.Bool = true

	validEmailVerificationRequestMockReturn := models.EmailVerificationRequest{
		ID:     1,
		UserID: 1,
		Email:  email,
		Code:   "123456",
		Expiry: pgtype.Timestamptz{
			Time:  time.Now().Add(time.Hour),
			Valid: true,
		},
	}

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
			email:     email,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				queries.EXPECT().GetUserByEmail(t.Context(), email).
					Return(validUserMockReturn, nil).Once()
				queries.EXPECT().InvalidateEmailVerificationRequestsOfUser(t.Context(), mock.Anything).
					Return(nil).Once()
				queries.EXPECT().GetEmailVerificationRequestByCode(t.Context(), mock.Anything).
					Return(models.EmailVerificationRequest{}, pgx.ErrNoRows).Once()
				queries.EXPECT().CreateEmailVerificationRequest(t.Context(), mock.Anything).
					Return(validEmailVerificationRequestMockReturn, nil).Once()
			},
		},
		{
			name:      "No user with email",
			email:     email,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				queries.EXPECT().GetUserByEmail(t.Context(), email).
					Return(models.User{}, pgx.ErrNoRows).Once()
			},
		},
		{
			name:      "Email already verified",
			email:     email,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				queries.EXPECT().GetUserByEmail(t.Context(), email).
					Return(emailVerifiedUserMockReturn, pgx.ErrNoRows).Once()
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
			email:     email,
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
