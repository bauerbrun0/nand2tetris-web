package userhandlers_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers/userhandlers"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleUserSettingsDeleteAccountPost(t *testing.T) {
	ts, queries, githubOauthService, googleOauthService := testutils.NewTestServer(t,
		testutils.TestServerOptions{
			Logs: false,
		},
	)
	defer ts.Close()

	_, csrfToken := ts.MustLogIn(t, testutils.LoginParams{})

	tests := []struct {
		name         string
		email        string
		verification string
		password     string
		csrfToken    string
		wantCode     int
		before       func(t *testing.T)
		after        func(t *testing.T)
	}{
		{
			name:         "With password verification",
			email:        testutils.MockEmail,
			verification: string(userhandlers.VerificationPassword),
			password:     testutils.MockPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
			before: func(t *testing.T) {
				testutils.ExpectGetUserByIdReturnsUser(t, queries)
				queries.EXPECT().DeleteUser(t.Context(), testutils.MockUserId).
					Return(nil).Once()
			},
			after: func(t *testing.T) {
				ts.MustLogIn(t, testutils.LoginParams{})
			},
		},
		{
			name:         "With GitHub verification",
			email:        testutils.MockEmail,
			verification: string(userhandlers.VerificationGitHub),
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
			before: func(t *testing.T) {
				githubOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, mock.Anything).
					Return(ts.URL).Once()
			},
			after: func(t *testing.T) {
				ts.MustLogIn(t, testutils.LoginParams{})
			},
		},
		{
			name:         "With Google verification",
			email:        testutils.MockEmail,
			verification: string(userhandlers.VerificationGoogle),
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
			before: func(t *testing.T) {
				googleOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, mock.Anything).
					Return(ts.URL).Once()
			},
			after: func(t *testing.T) {
				ts.MustLogIn(t, testutils.LoginParams{})
			},
		},
		{
			name:         "Empty email",
			email:        "",
			verification: string(userhandlers.VerificationPassword),
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
		},
		{
			name:         "Wrong email",
			email:        "wrong" + testutils.MockEmail,
			verification: string(userhandlers.VerificationPassword),
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
		},
		{
			name:         "Invalid verification",
			email:        testutils.MockEmail,
			verification: "invalid",
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}

			form := url.Values{}
			form.Add("Action", string(userhandlers.ActionDeleteAccount))
			form.Add("Verification", tt.verification)
			form.Add("csrf_token", tt.csrfToken)
			form.Add("DeleteAccount.Email", tt.email)
			form.Add("DeleteAccount.Password", tt.password)

			result := ts.PostForm(t, "/user/settings", form)
			assert.Equal(t, tt.wantCode, result.Status)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}

func TestUserDeleteAccountActionOAuthCallback(t *testing.T) {
	ts, queries, githubOauthService, googleOauthService := testutils.NewTestServer(t,
		testutils.TestServerOptions{
			Logs: false,
		},
	)
	defer ts.Close()

	result := ts.Get(t, "/user/login")
	csrfToken := testutils.ExtractCSRFToken(t, result.Body)

	var currentState string

	beforeEach := func(t *testing.T) {
		ts.MustLogIn(t, testutils.LoginParams{})
	}

	tests := []struct {
		name         string
		callbackPath string
		wantCode     int
		before       func(t *testing.T)
		after        func(t *testing.T)
	}{
		{
			name:         "Valid GitHub callback",
			callbackPath: "/user/oauth/github/callback/action",
			wantCode:     http.StatusSeeOther,
			before: func(t *testing.T) {
				// this sends a post request to /user/settings with Aciton = "delete-account"
				// and Verification = "GitHub" or "Google"
				// which responds with status See Other
				// function will return the state which should be used when sending a request
				// to /user/oauth/github|google/callback/action as a query parameter
				currentState = ts.MustSendUserSettingsOAuthAction(t, githubOauthService, googleOauthService, testutils.UserSettingsOAuthActionParams{
					Action:       userhandlers.ActionDeleteAccount,
					Verification: userhandlers.VerificationGitHub,
					CSRFToken:    csrfToken,
					FormData: map[string]string{
						"DeleteAccount.Email": testutils.MockEmail,
					},
				})

				testutils.ExpectExchangeCodeForUserInfo(t, githubOauthService)
				testutils.ExpectFindOAuthAuthorizationReturnsAuthorization(t, queries, models.ProviderGitHub)
				queries.EXPECT().DeleteUser(t.Context(), testutils.MockUserId).
					Return(nil).Once()
			},
		},
		{
			name:         "Valid Google callback",
			callbackPath: "/user/oauth/google/callback/action",
			wantCode:     http.StatusSeeOther,
			before: func(t *testing.T) {
				// this sends a post request to /user/settings with Aciton = "delete-account"
				// and Verification = "GitHub" or "Google"
				// which responds with status See Other
				// function will return the state which should be used when sending a request
				// to /user/oauth/github|google/callback/action as a query parameter
				currentState = ts.MustSendUserSettingsOAuthAction(t, githubOauthService, googleOauthService, testutils.UserSettingsOAuthActionParams{
					Action:       userhandlers.ActionDeleteAccount,
					Verification: userhandlers.VerificationGoogle,
					CSRFToken:    csrfToken,
					FormData: map[string]string{
						"DeleteAccount.Email": testutils.MockEmail,
					},
				})

				testutils.ExpectExchangeCodeForUserInfo(t, googleOauthService)
				testutils.ExpectFindOAuthAuthorizationReturnsAuthorization(t, queries, models.ProviderGoogle)
				queries.EXPECT().DeleteUser(t.Context(), testutils.MockUserId).
					Return(nil).Once()
			},
		},
		{
			name:         "OAuth authorization not exist",
			callbackPath: "/user/oauth/github/callback/action",
			wantCode:     http.StatusUnauthorized,
			before: func(t *testing.T) {
				// this sends a post request to /user/settings with Aciton = "delete-account"
				// and Verification = "GitHub" or "Google"
				// which responds with status See Other
				// function will return the state which should be used when sending a request
				// to /user/oauth/github|google/callback/action as a query parameter
				currentState = ts.MustSendUserSettingsOAuthAction(t, githubOauthService, googleOauthService, testutils.UserSettingsOAuthActionParams{
					Action:       userhandlers.ActionDeleteAccount,
					Verification: userhandlers.VerificationGitHub,
					CSRFToken:    csrfToken,
					FormData: map[string]string{
						"DeleteAccount.Email": testutils.MockEmail,
					},
				})

				testutils.ExpectExchangeCodeForUserInfo(t, githubOauthService)
				queries.EXPECT().FindOAuthAuthorization(t.Context(), models.FindOAuthAuthorizationParams{
					UserProviderID: testutils.MockOAuthUserId,
					Provider:       models.ProviderGitHub,
				}).Return(models.OauthAuthorization{}, pgx.ErrNoRows).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeEach(t)

			if tt.before != nil {
				tt.before(t)
			}

			path := fmt.Sprintf("%s?code=%s&state=%s", tt.callbackPath, testutils.MockOAuthCode, currentState)
			result = ts.Get(t, path)
			assert.Equal(t, tt.wantCode, result.Status)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}
