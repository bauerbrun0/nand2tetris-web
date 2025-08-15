package handlers_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers"
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

	_, csrfToken := ts.MustLogIn(t, testutils.LoginParams{})

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
			currentPassword:         testutils.MockPassword,
			newPassword:             testutils.MockPassword + "new",
			newPasswordConfirmation: testutils.MockPassword + "new",
			wantCode:                http.StatusSeeOther,
			before: func(t *testing.T) {
				queries.EXPECT().GetUserById(t.Context(), testutils.MockUserId).
					Return(models.User{
						ID:       testutils.MockUserId,
						Username: testutils.MockUsername,
						Email:    testutils.MockEmail,
						EmailVerified: pgtype.Bool{
							Bool:  true,
							Valid: true,
						},
						PasswordHash: pgtype.Text{
							String: testutils.MustHashPassword(t, testutils.MockPassword),
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
				ts.MustLogIn(t, testutils.LoginParams{})
			},
		},
		{
			name:                    "Wrong current password",
			currentPassword:         testutils.MockPassword + "wrong",
			newPassword:             testutils.MockPassword + "new",
			newPasswordConfirmation: testutils.MockPassword + "new",
			wantCode:                http.StatusUnauthorized,
			before: func(t *testing.T) {
				queries.EXPECT().GetUserById(t.Context(), testutils.MockUserId).
					Return(models.User{
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
							Time:  time.Now().Add(-time.Minute),
							Valid: true,
						},
					}, nil).Once()
			},
		},
		{
			name:                    "Password not set",
			currentPassword:         testutils.MockPassword,
			newPassword:             testutils.MockPassword + "new",
			newPasswordConfirmation: testutils.MockPassword + "new",
			wantCode:                http.StatusInternalServerError,
			before: func(t *testing.T) {
				queries.EXPECT().GetUserById(t.Context(), testutils.MockUserId).
					Return(models.User{
						ID:       testutils.MockUserId,
						Username: testutils.MockUsername,
						Email:    testutils.MockEmail,
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
				ts.MustLogIn(t, testutils.LoginParams{})
			},
		},
		{
			name:                    "Empty current password",
			currentPassword:         "",
			newPassword:             testutils.MockPassword + "new",
			newPasswordConfirmation: testutils.MockPassword + "new",
			wantCode:                http.StatusUnprocessableEntity,
		},
		{
			name:                    "Empty new password",
			currentPassword:         testutils.MockPassword,
			newPassword:             "",
			newPasswordConfirmation: testutils.MockPassword + "new",
			wantCode:                http.StatusUnprocessableEntity,
		},
		{
			name:                    "Empty new password confirmation",
			currentPassword:         testutils.MockPassword,
			newPassword:             testutils.MockPassword + "new",
			newPasswordConfirmation: "",
			wantCode:                http.StatusUnprocessableEntity,
		},
		{
			name:                    "New passwords do not match",
			currentPassword:         testutils.MockPassword,
			newPassword:             testutils.MockPassword + "new",
			newPasswordConfirmation: testutils.MockPassword + "new2",
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

			result := ts.PostForm(t, "/user/settings", form)
			assert.Equal(t, tt.wantCode, result.Status)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}
