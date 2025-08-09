package handlers_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers"
	"github.com/bauerbrun0/nand2tetris-web/internal/crypto"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleUserSettingsChangePasswordPost(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
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
		name                    string
		currentPassword         string
		newPassword             string
		newPasswordConfirmation string
		wantCode                int
		before                  func(t *testing.T)
		after                   func(t *testing.T)
	}{
		{
			name:                    "Valid submission",
			currentPassword:         password,
			newPassword:             password + "new",
			newPasswordConfirmation: password + "new",
			wantCode:                http.StatusSeeOther,
			before: func(t *testing.T) {
				var hasher crypto.PasswordHasher
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
							String: testutils.MustHashPassword(t, hasher, password),
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
			name:                    "Wrong current password",
			currentPassword:         password + "wrong",
			newPassword:             password + "new",
			newPasswordConfirmation: password + "new",
			wantCode:                http.StatusUnauthorized,
			before: func(t *testing.T) {
				var hasher crypto.PasswordHasher
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
							String: testutils.MustHashPassword(t, hasher, password),
							Valid:  true,
						},
						Created: pgtype.Timestamptz{
							Time:  time.Now().Add(-time.Minute),
							Valid: true,
						},
					}, nil).Once()
			},
		},
		{
			name:                    "Password not set",
			currentPassword:         password,
			newPassword:             password + "new",
			newPasswordConfirmation: password + "new",
			wantCode:                http.StatusInternalServerError,
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
			name:                    "Empty current password",
			currentPassword:         "",
			newPassword:             password + "new",
			newPasswordConfirmation: password + "new",
			wantCode:                http.StatusUnprocessableEntity,
		},
		{
			name:                    "Empty new password",
			currentPassword:         password,
			newPassword:             "",
			newPasswordConfirmation: password + "new",
			wantCode:                http.StatusUnprocessableEntity,
		},
		{
			name:                    "Empty new password confirmation",
			currentPassword:         password,
			newPassword:             password + "new",
			newPasswordConfirmation: "",
			wantCode:                http.StatusUnprocessableEntity,
		},
		{
			name:                    "New passwords do not match",
			currentPassword:         password,
			newPassword:             password + "new",
			newPasswordConfirmation: password + "new2",
			wantCode:                http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}

			form := url.Values{}
			form.Add("Action", string(handlers.ActionChangePassword))
			form.Add("csrf_token", csrfToken)
			form.Add("ChangePassword.CurrentPassword", tt.currentPassword)
			form.Add("ChangePassword.NewPassword", tt.newPassword)
			form.Add("ChangePassword.NewPasswordConfirmation", tt.newPasswordConfirmation)

			code, _, _ := ts.PostForm(t, "/user/settings", form)
			assert.Equal(t, tt.wantCode, code)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}
