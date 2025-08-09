package handlers_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	modelsmocks "github.com/bauerbrun0/nand2tetris-web/internal/models/mocks"
	servicemocks "github.com/bauerbrun0/nand2tetris-web/internal/services/mocks"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleUserSettingsCreatePasswordPost(t *testing.T) {
	githubOauthService := servicemocks.NewMockOAuthService(t)
	googleOauthService := servicemocks.NewMockOAuthService(t)
	queries := modelsmocks.NewMockDBQueries(t)

	ts := testutils.NewTestServer(t, queries, githubOauthService, googleOauthService, false)
	defer ts.Close()

	var (
		username = "walter"
		email    = "walter.white@example.com"
		password = "LosPollos321"
	)
	ts.MustLogIn(t, queries, testutils.LoginUser{
		Username: username,
		Email:    email,
		Password: password,
	})
	code, _, body := ts.Get(t, "/user/settings")
	assert.Equal(t, http.StatusOK, code)
	csrfToken := testutils.ExtractCSRFToken(t, body)
	assert.NotEmpty(t, csrfToken)

	tests := []struct {
		name                 string
		password             string
		passwordConfirmation string
		csrfToken            string
		wantCode             int
		before               func(t *testing.T)
		after                func(t *testing.T)
	}{
		{
			name:                 "Valid submission",
			password:             password,
			passwordConfirmation: password,
			csrfToken:            csrfToken,
			wantCode:             http.StatusSeeOther,
			before: func(t *testing.T) {
				queries.EXPECT().GetUserById(t.Context(), int32(1)).
					Return(models.User{
						ID:       1,
						Username: username,
						Email:    email,
						EmailVerified: pgtype.Bool{
							Bool:  true,
							Valid: true,
						},
						PasswordHash: pgtype.Text{
							String: "",
							Valid:  true,
						},
						Created: pgtype.Timestamptz{
							Time:  time.Now().Add(-time.Minute),
							Valid: true,
						},
					}, nil).Once()
				queries.EXPECT().ChangeUserPasswordHash(t.Context(), mock.Anything).
					Return(nil).Once()
			},
			after: func(t *testing.T) {
				ts.MustLogIn(t, queries, testutils.LoginUser{
					Username: username,
					Email:    email,
					Password: password,
				})
			},
		},
		{
			name:                 "Password already set",
			password:             password,
			passwordConfirmation: password,
			csrfToken:            csrfToken,
			wantCode:             http.StatusInternalServerError,
			before: func(t *testing.T) {
				queries.EXPECT().GetUserById(t.Context(), int32(1)).
					Return(models.User{
						ID:       1,
						Username: username,
						Email:    email,
						EmailVerified: pgtype.Bool{
							Bool:  true,
							Valid: true,
						},
						PasswordHash: pgtype.Text{
							String: "passwordhash",
							Valid:  true,
						},
						Created: pgtype.Timestamptz{
							Time:  time.Now().Add(-time.Minute),
							Valid: true,
						},
					}, nil).Once()
			},
			after: func(t *testing.T) {
				ts.MustLogIn(t, queries, testutils.LoginUser{
					Username: username,
					Email:    email,
					Password: password,
				})
			},
		},
		{
			name:                 "Empty password",
			password:             "",
			passwordConfirmation: password,
			csrfToken:            csrfToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Empty password confirmation",
			password:             password,
			passwordConfirmation: "",
			csrfToken:            csrfToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Passwords do not match",
			password:             password,
			passwordConfirmation: password + "x",
			csrfToken:            csrfToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Too long password",
			password:             strings.Repeat("x", 100),
			passwordConfirmation: strings.Repeat("x", 100),
			csrfToken:            csrfToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}

			form := url.Values{}
			form.Add("Action", string(handlers.ActionCreatePassword))
			form.Add("csrf_token", tt.csrfToken)
			form.Add("CreatePassword.Password", tt.password)
			form.Add("CreatePassword.PasswordConfirmation", tt.passwordConfirmation)

			code, _, _ := ts.PostForm(t, "/user/settings", form)
			assert.Equal(t, tt.wantCode, code)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}
