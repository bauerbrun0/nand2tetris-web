package userhandlers_test

import (
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

func TestHandleUserSettingsLinkAccountPost(t *testing.T) {
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
		action       userhandlers.Action
		wantCode     int
		before       func(t *testing.T)
		after        func(t *testing.T)
	}{
		{
			name:         "GitHub - with password",
			verification: string(userhandlers.VerificationPassword),
			action:       userhandlers.ActionLinkGitHubAccount,
			password:     testutils.MockPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
			before: func(t *testing.T) {
				testutils.ExpectGetUserByIdReturnsUser(t, queries)
				githubOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, "/user/oauth/github/callback/link").
					Return("https://github.com/oauth").Once()
			},
		},
		{
			name:         "GitHub - with wrong password",
			verification: string(userhandlers.VerificationPassword),
			action:       userhandlers.ActionLinkGitHubAccount,
			password:     testutils.MockPassword + "wrong",
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnauthorized,
			before: func(t *testing.T) {
				testutils.ExpectGetUserByIdReturnsUser(t, queries)
			},
		},
		{
			name:         "GitHub - with google",
			verification: string(userhandlers.VerificationGoogle),
			action:       userhandlers.ActionLinkGitHubAccount,
			password:     testutils.MockPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
			before: func(t *testing.T) {
				googleOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, "/user/oauth/google/callback/action").
					Return("https://google.com/oauth").Once()
			},
		},
		{
			name:         "GitHub - with github",
			verification: string(userhandlers.VerificationGitHub),
			action:       userhandlers.ActionLinkGitHubAccount,
			password:     testutils.MockPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Google - with password",
			verification: string(userhandlers.VerificationPassword),
			action:       userhandlers.ActionLinkGoogleAccount,
			password:     testutils.MockPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
			before: func(t *testing.T) {
				testutils.ExpectGetUserByIdReturnsUser(t, queries)
				googleOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, "/user/oauth/google/callback/link").
					Return("https://google.com/oauth").Once()
			},
		},
		{
			name:         "Google - with wrong password",
			verification: string(userhandlers.VerificationPassword),
			action:       userhandlers.ActionLinkGoogleAccount,
			password:     testutils.MockPassword + "wrong",
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnauthorized,
			before: func(t *testing.T) {
				testutils.ExpectGetUserByIdReturnsUser(t, queries)
			},
		},
		{
			name:         "Google - with github",
			verification: string(userhandlers.VerificationGitHub),
			action:       userhandlers.ActionLinkGoogleAccount,
			password:     testutils.MockPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
			before: func(t *testing.T) {
				githubOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, "/user/oauth/github/callback/action").
					Return("https://github.com/oauth").Once()
			},
		},
		{
			name:         "Google - with google",
			verification: string(userhandlers.VerificationGoogle),
			action:       userhandlers.ActionLinkGoogleAccount,
			password:     testutils.MockPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}

			form := url.Values{}
			form.Add("Action", string(tt.action))
			form.Add("Verification", tt.verification)
			form.Add("csrf_token", tt.csrfToken)

			switch tt.action {
			case userhandlers.ActionLinkGitHubAccount:
				form.Add("LinkGithub.Password", tt.password)
			case userhandlers.ActionLinkGoogleAccount:
				form.Add("LinkGoogle.Password", tt.password)
			default:
				t.Fatal("Invalid action")
			}

			result := ts.PostForm(t, "/user/settings", form)
			assert.Equal(t, tt.wantCode, result.Status)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}

func TestUserLinkOAuthCallback(t *testing.T) {
	ts, queries, githubOauthService, googleOauthService := testutils.NewTestServer(t,
		testutils.TestServerOptions{
			Logs: false,
		},
	)
	defer ts.Close()

	result := ts.Get(t, "/user/login")
	assert.Equal(t, http.StatusOK, result.Status)
	csrfToken := testutils.ExtractCSRFToken(t, result.Body)
	assert.NotEmpty(t, csrfToken)

	var currentState string

	ts.MustLogIn(t, testutils.LoginParams{})

	tests := []struct {
		name             string
		callbackPath     string
		wantCode         int
		wantRedirectPath string
		before           func(t *testing.T)
		after            func(t *testing.T)
	}{
		{
			name:             "Link GitHub account using Google verification",
			callbackPath:     "/user/oauth/github/callback/link",
			wantCode:         http.StatusSeeOther,
			wantRedirectPath: "/user/settings",
			before: func(t *testing.T) {
				currentState = ts.MustSendUserSettingsOAuthAction(t, githubOauthService, googleOauthService, testutils.UserSettingsOAuthActionParams{
					Action:       userhandlers.ActionLinkGitHubAccount,
					Verification: userhandlers.VerificationGoogle,
					CSRFToken:    csrfToken,
				})

				ts.MustAuthenticateOAuthAction(t, testutils.AuthenticateOAuthActionParams{
					State:        currentState,
					Verification: userhandlers.VerificationGoogle,
					BeforeActionRedirect: func() {
						githubOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, mock.Anything).
							RunAndReturn(func(state string, callbackPath string) string {
								currentState = state
								return "https://github.com/oauth"
							}).Once()
					},
				})

				testutils.ExpectExchangeCodeForUserInfo(t, githubOauthService)
				queries.EXPECT().FindOAuthAuthorization(t.Context(), models.FindOAuthAuthorizationParams{
					UserProviderID: testutils.MockOAuthUserId,
					Provider:       models.ProviderGitHub,
				}).Return(models.OauthAuthorization{}, pgx.ErrNoRows).Once()

				testutils.ExpectCreateOAuthAuthorizationReturnsAuthorization(t, queries, models.ProviderGitHub)
			},
		},
		{
			name:             "Link Google account using GitHub verification",
			callbackPath:     "/user/oauth/google/callback/link",
			wantCode:         http.StatusSeeOther,
			wantRedirectPath: "/user/settings",
			before: func(t *testing.T) {
				currentState = ts.MustSendUserSettingsOAuthAction(t, githubOauthService, googleOauthService, testutils.UserSettingsOAuthActionParams{
					Action:       userhandlers.ActionLinkGoogleAccount,
					Verification: userhandlers.VerificationGitHub,
					CSRFToken:    csrfToken,
				})

				ts.MustAuthenticateOAuthAction(t, testutils.AuthenticateOAuthActionParams{
					State:        currentState,
					Verification: userhandlers.VerificationGitHub,
					BeforeActionRedirect: func() {
						googleOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, mock.Anything).
							RunAndReturn(func(state string, callbackPath string) string {
								currentState = state
								return "https://google.com/oauth"
							}).Once()
					},
				})

				testutils.ExpectExchangeCodeForUserInfo(t, googleOauthService)

				queries.EXPECT().FindOAuthAuthorization(t.Context(), models.FindOAuthAuthorizationParams{
					UserProviderID: testutils.MockOAuthUserId,
					Provider:       models.ProviderGoogle,
				}).Return(models.OauthAuthorization{}, pgx.ErrNoRows).Once()

				testutils.ExpectCreateOAuthAuthorizationReturnsAuthorization(t, queries, models.ProviderGoogle)
			},
		},
		{
			name:             "GitHub account is already linked to an account",
			callbackPath:     "/user/oauth/github/callback/link",
			wantCode:         http.StatusSeeOther,
			wantRedirectPath: "/user/settings",
			before: func(t *testing.T) {
				currentState = ts.MustSendUserSettingsOAuthAction(t, githubOauthService, googleOauthService, testutils.UserSettingsOAuthActionParams{
					Action:       userhandlers.ActionLinkGitHubAccount,
					Verification: userhandlers.VerificationGoogle,
					CSRFToken:    csrfToken,
				})

				ts.MustAuthenticateOAuthAction(t, testutils.AuthenticateOAuthActionParams{
					State:        currentState,
					Verification: userhandlers.VerificationGoogle,
					BeforeActionRedirect: func() {
						githubOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, mock.Anything).
							RunAndReturn(func(state string, callbackPath string) string {
								currentState = state
								return "https://github.com/oauth"
							}).Once()
					},
				})

				testutils.ExpectExchangeCodeForUserInfo(t, githubOauthService)
				testutils.ExpectFindOAuthAuthorizationReturnsAuthorization(t, queries, models.ProviderGitHub)
			},
		},
		{
			name:             "Google account is already linked to an account",
			callbackPath:     "/user/oauth/google/callback/link",
			wantCode:         http.StatusSeeOther,
			wantRedirectPath: "/user/settings",
			before: func(t *testing.T) {
				currentState = ts.MustSendUserSettingsOAuthAction(t, githubOauthService, googleOauthService, testutils.UserSettingsOAuthActionParams{
					Action:       userhandlers.ActionLinkGoogleAccount,
					Verification: userhandlers.VerificationGitHub,
					CSRFToken:    csrfToken,
				})

				ts.MustAuthenticateOAuthAction(t, testutils.AuthenticateOAuthActionParams{
					State:        currentState,
					Verification: userhandlers.VerificationGitHub,
					BeforeActionRedirect: func() {
						googleOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, mock.Anything).
							RunAndReturn(func(state string, callbackPath string) string {
								currentState = state
								return "https://google.com/oauth"
							}).Once()
					},
				})

				testutils.ExpectExchangeCodeForUserInfo(t, googleOauthService)
				testutils.ExpectFindOAuthAuthorizationReturnsAuthorization(t, queries, models.ProviderGoogle)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}

			res := ts.Get(t, tt.callbackPath+"?state="+currentState+"&code="+testutils.MockOAuthCode)
			assert.Equal(t, http.StatusSeeOther, res.Status)
			if tt.wantRedirectPath != "" {
				assert.Equal(t, ts.URL+tt.wantRedirectPath, res.RedirectUrl.String())
			}
			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}
