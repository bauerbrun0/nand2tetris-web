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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserRegister(t *testing.T) {
	ts, _, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	t.Run("Can visit page", func(t *testing.T) {
		result := ts.Get(t, "/user/register")

		assert.Equal(t, http.StatusOK, result.Status, "status code should be 200 OK")

		csrfToken := testutils.ExtractCSRFToken(t, result.Body)
		assert.NotEmptyf(t, csrfToken, "csrfToken should not be empty")
	})

	t.Run("Redirect if already logged in", func(t *testing.T) {
		ts.MustLogIn(t, testutils.LoginParams{})
		result := ts.Get(t, "/user/register")
		assert.Equal(t, http.StatusSeeOther, result.Status, "status code should be 303 See Other")
	})
}

func TestUserRegisterPost(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	result := ts.Get(t, "/user/register")
	validCSRFToken := testutils.ExtractCSRFToken(t, result.Body)

	returnUser := models.User{
		ID:       testutils.MockUserId,
		Username: testutils.MockUsername,
		Email:    testutils.MockEmail,
		EmailVerified: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		PasswordHash: pgtype.Text{
			String: testutils.MockPasswordHash,
			Valid:  true,
		},
		Created: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	returnEmailVerificationRequest := models.EmailVerificationRequest{
		ID:     testutils.MockId,
		UserID: testutils.MockUserId,
		Email:  testutils.MockEmail,
		Code:   "123456",
		Expiry: pgtype.Timestamptz{
			Time:  time.Now().Add(time.Hour),
			Valid: true,
		},
	}

	tests := []struct {
		name                 string
		username             string
		email                string
		password             string
		passwordConfirmation string
		terms                string
		csrfToken            string
		wantCode             int
		before               func(t *testing.T)
		after                func(t *testing.T)
	}{
		{
			name:                 "Valid submission",
			username:             testutils.MockUsername,
			email:                testutils.MockEmail,
			password:             testutils.MockPassword,
			passwordConfirmation: testutils.MockPassword,
			terms:                testutils.MockTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusSeeOther,
			before: func(t *testing.T) {
				queries.EXPECT().CreateNewUser(t.Context(), mock.Anything).
					Return(returnUser, nil).Once()
				queries.EXPECT().GetEmailVerificationRequestByCode(t.Context(), mock.Anything).
					Return(models.EmailVerificationRequest{}, pgx.ErrNoRows)
				queries.EXPECT().CreateEmailVerificationRequest(t.Context(), mock.Anything).
					Return(returnEmailVerificationRequest, nil)
			},
		},
		{
			name:                 "Duplicate email",
			username:             testutils.MockUsername,
			email:                testutils.MockEmail,
			password:             testutils.MockPassword,
			passwordConfirmation: testutils.MockPassword,
			terms:                testutils.MockTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
			before: func(t *testing.T) {
				queries.EXPECT().CreateNewUser(t.Context(), mock.Anything).
					Return(
						models.User{},
						&pgconn.PgError{
							Code:           models.ErrorCodeUniqueViolation,
							ConstraintName: models.ConstraintNameUsersUniqueEmail,
						},
					).
					Once()
			},
		},
		{
			name:                 "Duplicate username",
			username:             testutils.MockUsername,
			email:                testutils.MockEmail,
			password:             testutils.MockPassword,
			passwordConfirmation: testutils.MockPassword,
			terms:                testutils.MockTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
			before: func(t *testing.T) {
				queries.EXPECT().CreateNewUser(t.Context(), mock.Anything).
					Return(
						models.User{},
						&pgconn.PgError{
							Code:           models.ErrorCodeUniqueViolation,
							ConstraintName: models.ConstraintNameUsersUniqueUsername,
						},
					).
					Once()
			},
		},
		{
			name:                 "Redirect if already logged in",
			username:             testutils.MockUsername,
			email:                testutils.MockEmail,
			password:             testutils.MockPassword,
			passwordConfirmation: testutils.MockPassword,
			terms:                testutils.MockTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusSeeOther,
			before: func(t *testing.T) {
				ts.MustLogIn(t, testutils.LoginParams{})
			},
			after: func(t *testing.T) {
				ts.RemoveCookie(t, "session")
			},
		},
		{
			name:                 "Empty password",
			username:             testutils.MockUsername,
			email:                testutils.MockEmail,
			password:             "",
			passwordConfirmation: testutils.MockPassword,
			terms:                testutils.MockTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Empty password confirmation",
			username:             testutils.MockUsername,
			email:                testutils.MockEmail,
			password:             testutils.MockPassword,
			passwordConfirmation: "",
			terms:                testutils.MockTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Passwords do not match",
			username:             testutils.MockUsername,
			email:                testutils.MockEmail,
			password:             testutils.MockPassword,
			passwordConfirmation: testutils.MockPassword + "extra",
			terms:                testutils.MockTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Password too short",
			username:             testutils.MockUsername,
			email:                testutils.MockEmail,
			password:             testutils.MockPassword,
			passwordConfirmation: "123",
			terms:                testutils.MockTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Empty username",
			username:             "",
			email:                testutils.MockEmail,
			password:             testutils.MockPassword,
			passwordConfirmation: testutils.MockPassword,
			terms:                testutils.MockTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Username too short",
			username:             "xy",
			email:                testutils.MockEmail,
			password:             testutils.MockPassword,
			passwordConfirmation: testutils.MockPassword,
			terms:                testutils.MockTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Username too long",
			username:             strings.Repeat("x", 100),
			email:                testutils.MockEmail,
			password:             testutils.MockPassword,
			passwordConfirmation: testutils.MockPassword,
			terms:                testutils.MockTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Username contains whitespace",
			username:             "walter white",
			email:                testutils.MockEmail,
			password:             testutils.MockPassword,
			passwordConfirmation: testutils.MockPassword,
			terms:                testutils.MockTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Empty email",
			username:             testutils.MockUsername,
			email:                "",
			password:             testutils.MockPassword,
			passwordConfirmation: testutils.MockPassword,
			terms:                testutils.MockTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Invalid email",
			username:             testutils.MockUsername,
			email:                "notemail@example",
			password:             testutils.MockPassword,
			passwordConfirmation: testutils.MockPassword,
			terms:                testutils.MockTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Terms not accepted",
			username:             testutils.MockUsername,
			email:                testutils.MockEmail,
			password:             testutils.MockPassword,
			passwordConfirmation: testutils.MockPassword,
			terms:                "",
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}

			form := url.Values{}
			form.Add("username", tt.username)
			form.Add("email", tt.email)
			form.Add("password", tt.password)
			form.Add("password-confirmation", tt.passwordConfirmation)
			form.Add("terms", tt.terms)
			form.Add("csrf_token", tt.csrfToken)

			result := ts.PostForm(t, "/user/register", form)
			assert.Equal(t, tt.wantCode, result.Status)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}

}
