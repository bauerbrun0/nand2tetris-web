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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserRegister(t *testing.T) {
	githubOauthService := servicemocks.NewMockOAuthService(t)
	googleOauthService := servicemocks.NewMockOAuthService(t)
	queries := modelsmocks.NewMockDBQueries(t)

	ts := testutils.NewTestServer(t, queries, githubOauthService, googleOauthService, false)
	defer ts.Close()

	t.Run("Can visit page", func(t *testing.T) {
		code, _, body := ts.Get(t, "/user/register")

		assert.Equal(t, http.StatusOK, code, "status code should be 200 OK")

		csrfToken := testutils.ExtractCSRFToken(t, body)
		assert.NotEmptyf(t, csrfToken, "csrfToken should not be empty")
	})

	t.Run("Redirect if already logged in", func(t *testing.T) {
		ts.MustLogIn(t, queries, "walter", "walter.white@example.com", "LosPollos321")
		code, _, _ := ts.Get(t, "/user/register")
		assert.Equal(t, http.StatusSeeOther, code, "status code should be 303 See Other")
	})
}

func TestUserRegisterPost(t *testing.T) {
	githubOauthService := servicemocks.NewMockOAuthService(t)
	googleOauthService := servicemocks.NewMockOAuthService(t)
	queries := modelsmocks.NewMockDBQueries(t)

	ts := testutils.NewTestServer(t, queries, githubOauthService, googleOauthService, false)
	defer ts.Close()

	_, _, body := ts.Get(t, "/user/register")
	validCSRFToken := testutils.ExtractCSRFToken(t, body)

	var hasher crypto.PasswordHasher

	var (
		validUsername     = "walter"
		validEmail        = "walter.white@example.com"
		validPassword     = "LosPollos321"
		validTerms        = "on"
		validPasswordHash = testutils.MustHashPassword(t, hasher, validPassword)
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

	returnEmailVerificationRequest := models.EmailVerificationRequest{
		ID:     1,
		UserID: 1,
		Email:  validEmail,
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
			username:             validUsername,
			email:                validEmail,
			password:             validPassword,
			passwordConfirmation: validPassword,
			terms:                validTerms,
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
			username:             validUsername,
			email:                validEmail,
			password:             validPassword,
			passwordConfirmation: validPassword,
			terms:                validTerms,
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
			username:             validUsername,
			email:                validEmail,
			password:             validPassword,
			passwordConfirmation: validPassword,
			terms:                validTerms,
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
			username:             validUsername,
			email:                validEmail,
			password:             validPassword,
			passwordConfirmation: validPassword,
			terms:                validTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusSeeOther,
			before: func(t *testing.T) {
				ts.MustLogIn(t, queries, validUsername, validEmail, validPassword)
			},
			after: func(t *testing.T) {
				ts.RemoveCookie(t, "session")
			},
		},
		{
			name:                 "Empty password",
			username:             validUsername,
			email:                validEmail,
			password:             "",
			passwordConfirmation: validPassword,
			terms:                validTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Empty password confirmation",
			username:             validUsername,
			email:                validEmail,
			password:             validPassword,
			passwordConfirmation: "",
			terms:                validTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Passwords do not match",
			username:             validUsername,
			email:                validEmail,
			password:             validPassword,
			passwordConfirmation: validPassword + "extra",
			terms:                validTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Password too short",
			username:             validUsername,
			email:                validEmail,
			password:             validPassword,
			passwordConfirmation: "123",
			terms:                validTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Empty username",
			username:             "",
			email:                validEmail,
			password:             validPassword,
			passwordConfirmation: validPassword,
			terms:                validTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Username too short",
			username:             "xy",
			email:                validEmail,
			password:             validPassword,
			passwordConfirmation: validPassword,
			terms:                validTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Username too long",
			username:             strings.Repeat("x", 100),
			email:                validEmail,
			password:             validPassword,
			passwordConfirmation: validPassword,
			terms:                validTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Username contains whitespace",
			username:             "walter white",
			email:                validEmail,
			password:             validPassword,
			passwordConfirmation: validPassword,
			terms:                validTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Empty email",
			username:             validUsername,
			email:                "",
			password:             validPassword,
			passwordConfirmation: validPassword,
			terms:                validTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Invalid email",
			username:             validUsername,
			email:                "notemail@example",
			password:             validPassword,
			passwordConfirmation: validPassword,
			terms:                validTerms,
			csrfToken:            validCSRFToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Terms not accepted",
			username:             validUsername,
			email:                validEmail,
			password:             validPassword,
			passwordConfirmation: validPassword,
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

			code, _, _ := ts.PostForm(t, "/user/register", form)
			assert.Equal(t, tt.wantCode, code)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}

}
