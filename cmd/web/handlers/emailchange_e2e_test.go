package handlers_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleUserSettingsChangeEmailPost(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	var (
		username = "walter"
		email    = "walter.white@example.com"
		newEmail = "walter.white@new.com"
		password = "LosPollos321"
	)
	ts.MustLogIn(t, queries, testutils.LoginUser{
		Username: username,
		Email:    email,
		Password: password,
	})
	result := ts.Get(t, "/user/settings")
	assert.Equal(t, http.StatusOK, result.Status)
	csrfToken := testutils.ExtractCSRFToken(t, result.Body)
	assert.NotEmpty(t, csrfToken)

	tests := []struct {
		name     string
		password string
		newEmail string
		wantCode int
		before   func(t *testing.T)
		after    func(t *testing.T)
	}{
		{
			name:     "Valid submission",
			password: password,
			newEmail: newEmail,
			wantCode: http.StatusOK,
			before: func(t *testing.T) {
				queries.EXPECT().GetUserById(t.Context(), int32(1)).Return(models.User{
					ID:       1,
					Username: username,
					Email:    email,
					EmailVerified: pgtype.Bool{
						Bool:  true,
						Valid: true,
					},
					PasswordHash: pgtype.Text{
						String: testutils.MustHashPassword(t, password),
						Valid:  true,
					},
					Created: pgtype.Timestamptz{
						Time:  time.Now().Add(-time.Hour),
						Valid: true,
					},
				}, nil).Once()
				queries.EXPECT().InvalidateEmailVerificationRequestsOfUser(t.Context(), mock.Anything).
					Return(nil).Once()
				queries.EXPECT().GetEmailVerificationRequestByCode(t.Context(), mock.Anything).
					Return(models.EmailVerificationRequest{}, pgx.ErrNoRows).Once()
				queries.EXPECT().CreateEmailVerificationRequest(t.Context(), mock.Anything).
					Return(models.EmailVerificationRequest{
						ID:     1,
						UserID: 1,
						Email:  email,
						Code:   "12345678",
						Expiry: pgtype.Timestamptz{
							Time:  time.Now().Add(time.Hour),
							Valid: true,
						}}, nil).
					Once()
			},
		},
		{
			name:     "Wrong password",
			password: password + "wrong",
			newEmail: newEmail,
			wantCode: http.StatusUnauthorized,
			before: func(t *testing.T) {
				queries.EXPECT().GetUserById(t.Context(), int32(1)).Return(models.User{
					ID:       1,
					Username: username,
					Email:    email,
					EmailVerified: pgtype.Bool{
						Bool:  true,
						Valid: true,
					},
					PasswordHash: pgtype.Text{
						String: testutils.MustHashPassword(t, password),
						Valid:  true,
					},
					Created: pgtype.Timestamptz{
						Time:  time.Now().Add(-time.Hour),
						Valid: true,
					},
				}, nil).Once()
			},
		},
		{
			name:     "Provide current email password",
			password: password + "wrong",
			newEmail: email,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "Empty password",
			password: "",
			newEmail: newEmail,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "Empty new email",
			password: password,
			newEmail: "",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "Invalid email",
			password: password,
			newEmail: email + "@",
			wantCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}

			form := url.Values{}
			form.Add("Action", string(handlers.ActionChangeEmail))
			form.Add("csrf_token", csrfToken)
			form.Add("ChangeEmail.Password", tt.password)
			form.Add("ChangeEmail.NewEmail", tt.newEmail)

			result := ts.PostForm(t, "/user/settings", form)
			assert.Equal(t, tt.wantCode, result.Status)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}

func TestHandleUserSettingsChangeEmailSendCodePost(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	var (
		username        = "walter"
		email           = "walter.white@example.com"
		password        = "LosPollos321"
		emailChangeCode = "12345678"
	)
	ts.MustLogIn(t, queries, testutils.LoginUser{
		Username: username,
		Email:    email,
		Password: password,
	})
	result := ts.Get(t, "/user/settings")
	assert.Equal(t, http.StatusOK, result.Status)
	csrfToken := testutils.ExtractCSRFToken(t, result.Body)
	assert.NotEmpty(t, csrfToken)

	tests := []struct {
		name     string
		code     string
		wantCode int
		before   func(t *testing.T)
		after    func(t *testing.T)
	}{
		{
			name:     "Valid submission",
			code:     emailChangeCode,
			wantCode: http.StatusSeeOther,
			before: func(t *testing.T) {
				queries.EXPECT().GetEmailVerificationRequestByCode(t.Context(), emailChangeCode).
					Return(models.EmailVerificationRequest{
						ID:     1,
						UserID: 1,
						Email:  email,
						Code:   emailChangeCode,
						Expiry: pgtype.Timestamptz{
							Time:  time.Now().Add(time.Hour),
							Valid: true,
						},
					}, nil).Once()
				queries.EXPECT().GetUserById(t.Context(), int32(1)).Return(models.User{
					ID:       1,
					Username: username,
					Email:    "user.oldemail@example.com",
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
				}, nil).Once()
				queries.EXPECT().InvalidateEmailVerificationRequest(t.Context(), mock.Anything).
					Return(nil).Once()
				queries.EXPECT().ChangeUserEmail(t.Context(), mock.Anything).
					Return(nil).Once()
			},
		},
		{
			name:     "Code expired",
			code:     emailChangeCode,
			wantCode: http.StatusUnauthorized,
			before: func(t *testing.T) {
				queries.EXPECT().GetEmailVerificationRequestByCode(t.Context(), emailChangeCode).
					Return(models.EmailVerificationRequest{
						ID:     1,
						UserID: 1,
						Email:  email,
						Code:   emailChangeCode,
						Expiry: pgtype.Timestamptz{
							Time:  time.Now().Add(-time.Hour),
							Valid: true,
						},
					}, nil).Once()
			},
		},
		{
			name:     "Code not exist",
			code:     emailChangeCode,
			wantCode: http.StatusUnauthorized,
			before: func(t *testing.T) {
				queries.EXPECT().GetEmailVerificationRequestByCode(t.Context(), emailChangeCode).
					Return(models.EmailVerificationRequest{}, pgx.ErrNoRows).Once()
			},
		},
		{
			name:     "Empty code",
			code:     "",
			wantCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}

			form := url.Values{}
			form.Add("Action", string(handlers.ActionChangeEmailSendCode))
			form.Add("csrf_token", csrfToken)
			form.Add("ChangeEmailSendCode.Code", tt.code)

			result := ts.PostForm(t, "/user/settings", form)
			assert.Equal(t, tt.wantCode, result.Status)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}
