package userhandlers_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers/userhandlers"
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

	_, csrfToken := ts.MustLogIn(t, testutils.LoginParams{})

	var newEmail = "walter.white@new.com"

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
			password: testutils.MockPassword,
			newEmail: newEmail,
			wantCode: http.StatusOK,
			before: func(t *testing.T) {
				testutils.ExpectGetUserByIdReturnsUser(t, queries)
				queries.EXPECT().InvalidateEmailVerificationRequestsOfUser(t.Context(), mock.Anything).
					Return(nil).Once()
				queries.EXPECT().GetEmailVerificationRequestByCode(t.Context(), mock.Anything).
					Return(models.EmailVerificationRequest{}, pgx.ErrNoRows).Once()
				testutils.ExpectCreateEmailVerificationRequestReturnsRequest(t, queries)
			},
		},
		{
			name:     "Wrong password",
			password: testutils.MockPassword + "wrong",
			newEmail: newEmail,
			wantCode: http.StatusUnauthorized,
			before: func(t *testing.T) {
				testutils.ExpectGetUserByIdReturnsUser(t, queries)
			},
		},
		{
			name:     "Provide current email password",
			password: testutils.MockPassword + "wrong",
			newEmail: testutils.MockEmail,
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
			password: testutils.MockPassword,
			newEmail: "",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "Invalid email",
			password: testutils.MockPassword,
			newEmail: testutils.MockEmail + "@",
			wantCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}

			form := url.Values{}
			form.Add("Action", string(userhandlers.ActionChangeEmail))
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

	_, csrfToken := ts.MustLogIn(t, testutils.LoginParams{})

	emailChangeCode := testutils.MockEmailVerificationRequestCode

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
				testutils.ExpectGetEmailVerificationRequestByCodeReturnsRequest(t, queries)
				testutils.ExpectGetUserByIdReturnsUser(t, queries)
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
						ID:     testutils.MockId,
						UserID: testutils.MockUserId,
						Email:  testutils.MockEmail,
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
			form.Add("Action", string(userhandlers.ActionChangeEmailSendCode))
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
