package userhandlers_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers/userhandlers"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleUserSettingsCreatePasswordPost(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	_, csrfToken := ts.MustLogIn(t, testutils.LoginParams{})

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
			password:             testutils.MockPassword,
			passwordConfirmation: testutils.MockPassword,
			csrfToken:            csrfToken,
			wantCode:             http.StatusSeeOther,
			before: func(t *testing.T) {
				testutils.ExpectGetUserByIdReturnsEmptyPasswordUser(t, queries)
				queries.EXPECT().ChangeUserPasswordHash(t.Context(), mock.Anything).
					Return(nil).Once()
			},
			after: func(t *testing.T) {
				ts.MustLogIn(t, testutils.LoginParams{})
			},
		},
		{
			name:                 "Password already set",
			password:             testutils.MockPassword,
			passwordConfirmation: testutils.MockPassword,
			csrfToken:            csrfToken,
			wantCode:             http.StatusInternalServerError,
			before: func(t *testing.T) {
				testutils.ExpectGetUserByIdReturnsUser(t, queries)
			},
		},
		{
			name:                 "Empty password",
			password:             "",
			passwordConfirmation: testutils.MockPassword,
			csrfToken:            csrfToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Empty password confirmation",
			password:             testutils.MockPassword,
			passwordConfirmation: "",
			csrfToken:            csrfToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Passwords do not match",
			password:             testutils.MockPassword,
			passwordConfirmation: testutils.MockPassword + "x",
			csrfToken:            csrfToken,
			wantCode:             http.StatusUnprocessableEntity,
		},
		{
			name:                 "Too long password",
			password:             testutils.MockLongPassword,
			passwordConfirmation: testutils.MockLongPassword,
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
			form.Add("Action", string(userhandlers.ActionCreatePassword))
			form.Add("csrf_token", tt.csrfToken)
			form.Add("CreatePassword.Password", tt.password)
			form.Add("CreatePassword.PasswordConfirmation", tt.passwordConfirmation)

			result := ts.PostForm(t, "/user/settings", form)
			assert.Equal(t, tt.wantCode, result.Status)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}
