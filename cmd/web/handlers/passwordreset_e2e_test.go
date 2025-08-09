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
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	t.Run("Can visit page", func(t *testing.T) {
		code, _, body := ts.Get(t, "/user/reset-password/send-code")
		assert.Equal(t, http.StatusOK, code)
		csrfToken := testutils.ExtractCSRFToken(t, body)
		assert.NotEmpty(t, csrfToken)
	})

	t.Run("Redirect if already logged in", func(t *testing.T) {
		var (
			username = "walt"
			email    = "walter.white@example.com"
			password = "LosPollos321"
		)
		ts.MustLogIn(t, queries, testutils.LoginUser{
			Username: username,
			Email:    email,
			Password: password,
		})
		code, _, _ := ts.Get(t, "/user/reset-password/send-code")
		assert.Equal(t, http.StatusSeeOther, code)
	})
}

func TestUserResetPasswordSendCodePost(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	_, _, body := ts.Get(t, "/user/reset-password/send-code")
	validCSRFToken := testutils.ExtractCSRFToken(t, body)

	var (
		username = "walter"
		email    = "walter.white@example.com"
		password = "LosPollos321"
	)

	userMockReturn := models.User{
		ID:       1,
		Username: username,
		Email:    email,
		EmailVerified: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		PasswordHash: pgtype.Text{
			String: "hash",
			Valid:  true,
		},
		Created: pgtype.Timestamptz{
			Time:  time.Now().Add(-time.Hour),
			Valid: true,
		},
	}
	passwordResetRequestMockReturn := models.PasswordResetRequest{
		ID:     1,
		UserID: 1,
		Email:  email,
		Code:   "123456789123",
		VerifyEmailAfter: pgtype.Bool{
			Bool:  false,
			Valid: true,
		},
		Expiry: pgtype.Timestamptz{
			Time:  time.Now().Add(time.Hour),
			Valid: true,
		},
	}

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
			email:     email,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				queries.EXPECT().GetUserByEmail(t.Context(), email).
					Return(userMockReturn, nil).Once()
				queries.EXPECT().InvalidatePasswordResetRequestsOfUser(t.Context(), mock.Anything).
					Return(nil).Once()
				queries.EXPECT().GetPasswordResetRequestByCode(t.Context(), mock.Anything).
					Return(models.PasswordResetRequest{}, pgx.ErrNoRows).Once()
				queries.EXPECT().CreatePasswordResetRequest(t.Context(), mock.Anything).
					Return(passwordResetRequestMockReturn, nil).Once()
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
			name:      "Redirect if already logged in",
			email:     email,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				ts.MustLogIn(t, queries, testutils.LoginUser{
					Username: username,
					Email:    email,
					Password: password,
				})
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

			code, _, _ := ts.PostForm(t, "/user/reset-password/send-code", form)
			assert.Equal(t, tt.wantCode, code)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}

func TestUserResetPasswordEnterCode(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	t.Run("Can visit page", func(t *testing.T) {
		code, _, body := ts.Get(t, "/user/reset-password/enter-code")
		assert.Equal(t, http.StatusOK, code)
		csrfToken := testutils.ExtractCSRFToken(t, body)
		assert.NotEmpty(t, csrfToken)
	})

	t.Run("Redirect if already logged in", func(t *testing.T) {
		var (
			username = "walt"
			email    = "walter.white@example.com"
			password = "LosPollos321"
		)
		ts.MustLogIn(t, queries, testutils.LoginUser{
			Username: username,
			Email:    email,
			Password: password,
		})
		code, _, _ := ts.Get(t, "/user/reset-password/enter-code")
		assert.Equal(t, http.StatusSeeOther, code)
	})
}

func TestUserResetPasswordEnterCodePost(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	_, _, body := ts.Get(t, "/user/reset-password/enter-code")
	validCSRFToken := testutils.ExtractCSRFToken(t, body)

	var (
		username = "walter"
		email    = "walter.white@example.com"
		code     = "12345678"
		password = "LosPollos321"
	)

	passwordResetRequestMockReturn := models.PasswordResetRequest{
		ID:     1,
		UserID: 1,
		Email:  email,
		Code:   code,
		VerifyEmailAfter: pgtype.Bool{
			Bool:  false,
			Valid: true,
		},
		Expiry: pgtype.Timestamptz{
			Time:  time.Now().Add(time.Hour),
			Valid: true,
		},
	}

	expiredPasswordResetRequestMockReturn := passwordResetRequestMockReturn
	expiredPasswordResetRequestMockReturn.Expiry.Time = time.Now().Add(-time.Minute)

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
			email:     email,
			code:      code,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				queries.EXPECT().GetPasswordResetRequestByCode(t.Context(), code).
					Return(passwordResetRequestMockReturn, nil).Once()
			},
		},
		{
			name:      "Redirect if already logged in",
			email:     email,
			code:      code,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				ts.MustLogIn(t, queries, testutils.LoginUser{
					Username: username,
					Email:    email,
					Password: password,
				})
			},
			after: func(t *testing.T) {
				ts.RemoveCookie(t, "session")
			},
		},
		{
			name:      "Code expired",
			email:     email,
			code:      code,
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
			before: func(t *testing.T) {
				queries.EXPECT().GetPasswordResetRequestByCode(t.Context(), code).
					Return(expiredPasswordResetRequestMockReturn, nil).Once()
			},
		},
		{
			name:      "Code does not exists",
			email:     email,
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
			email:     email,
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
			email:     email,
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

			code, _, _ := ts.PostForm(t, "/user/reset-password/enter-code", form)
			assert.Equal(t, tt.wantCode, code)

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
		code, _, _ := ts.Get(t, "/user/reset-password")
		assert.Equal(t, http.StatusSeeOther, code)
	})

	t.Run("Can visit page if code is in session", func(t *testing.T) {
		code, _, body := ts.Get(t, "/user/reset-password/enter-code")
		csrfToken := testutils.ExtractCSRFToken(t, body)

		var (
			email     = "walter.white@example.com"
			resetCode = "12345678"
		)

		form := url.Values{}
		form.Add("csrf_token", csrfToken)
		form.Add("email", email)
		form.Add("code", resetCode)

		passwordResetRequestMockReturn := models.PasswordResetRequest{
			ID:     1,
			UserID: 1,
			Email:  email,
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

		code, _, _ = ts.PostForm(t, "/user/reset-password/enter-code", form)
		assert.Equal(t, http.StatusSeeOther, code)

		code, _, _ = ts.Get(t, "/user/reset-password")
		assert.Equal(t, http.StatusOK, code)
	})

	t.Run("Redirect if already logged in", func(t *testing.T) {
		var (
			username = "walt"
			email    = "walter.white@example.com"
			password = "LosPollos321"
		)
		ts.MustLogIn(t, queries, testutils.LoginUser{
			Username: username,
			Email:    email,
			Password: password,
		})
		code, _, _ := ts.Get(t, "/user/reset-password")
		assert.Equal(t, http.StatusSeeOther, code)
	})
}

func TestUserResetPasswordPost(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	_, _, body := ts.Get(t, "/user/reset-password/enter-code")
	csrfToken := testutils.ExtractCSRFToken(t, body)

	var (
		resetCode = "12345678"
		username  = "walter"
		email     = "walter.white@example.com"
		password  = "LosPollos321"
	)

	form := url.Values{}
	form.Add("csrf_token", csrfToken)
	form.Add("email", email)
	form.Add("code", resetCode)

	passwordResetRequestMockReturn := models.PasswordResetRequest{
		ID:     1,
		UserID: 1,
		Email:  email,
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

	expiredPasswordResetRequestMockReturn := passwordResetRequestMockReturn
	expiredPasswordResetRequestMockReturn.Expiry.Time = time.Now().Add(-time.Minute)

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
			newPassword:             password,
			newPasswordConfirmation: password,
			csrfToken:               csrfToken,
			wantCode:                http.StatusSeeOther,
			before: func(t *testing.T) {
				queries.EXPECT().GetPasswordResetRequestByCode(t.Context(), resetCode).
					Return(passwordResetRequestMockReturn, nil).Once()
				queries.EXPECT().InvalidatePasswordResetRequest(t.Context(), mock.Anything).
					Return(nil).Once()
				queries.EXPECT().ChangeUserPasswordHash(t.Context(), mock.Anything).
					Return(nil).Once()
			},
		},
		{
			name:                    "Expired code",
			code:                    resetCode,
			newPassword:             password,
			newPasswordConfirmation: password,
			csrfToken:               csrfToken,
			wantCode:                http.StatusUnauthorized,
			before: func(t *testing.T) {
				queries.EXPECT().GetPasswordResetRequestByCode(t.Context(), resetCode).
					Return(expiredPasswordResetRequestMockReturn, nil).Once()
			},
		},
		{
			name:                    "Redirect if already logged in",
			code:                    resetCode,
			newPassword:             password,
			newPasswordConfirmation: password,
			csrfToken:               csrfToken,
			wantCode:                http.StatusSeeOther,
			before: func(t *testing.T) {
				ts.MustLogIn(t, queries, testutils.LoginUser{
					Username: username,
					Email:    email,
					Password: password,
				})
			},
			after: func(t *testing.T) {
				ts.RemoveCookie(t, "session")
			},
		},
		{
			name:                    "Empty code",
			code:                    "",
			newPassword:             password,
			newPasswordConfirmation: password,
			csrfToken:               csrfToken,
			wantCode:                http.StatusUnprocessableEntity,
		},
		{
			name:                    "Empty password",
			code:                    resetCode,
			newPassword:             "",
			newPasswordConfirmation: password,
			csrfToken:               csrfToken,
			wantCode:                http.StatusUnprocessableEntity,
		},
		{
			name:                    "Empty password confirmation",
			code:                    resetCode,
			newPassword:             password,
			newPasswordConfirmation: "",
			csrfToken:               csrfToken,
			wantCode:                http.StatusUnprocessableEntity,
		},
		{
			name:                    "Passwords do not match",
			code:                    resetCode,
			newPassword:             password,
			newPasswordConfirmation: password + "x",
			csrfToken:               csrfToken,
			wantCode:                http.StatusUnprocessableEntity,
		},
		{
			name:                    "Empty csrf token",
			code:                    resetCode,
			newPassword:             password,
			newPasswordConfirmation: password,
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

			code, _, _ := ts.PostForm(t, "/user/reset-password", form)
			assert.Equal(t, tt.wantCode, code)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}
