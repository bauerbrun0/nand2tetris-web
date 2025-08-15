package handlers_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleUserSettingsUnlinkAccountPost(t *testing.T) {
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
		action       handlers.Action
		wantCode     int
		before       func(t *testing.T)
		after        func(t *testing.T)
	}{
		{
			name:         "GitHub - with password",
			verification: string(handlers.VerificationPassword),
			action:       handlers.ActionUnlinkGitHubAccount,
			password:     testutils.MockPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
			before: func(t *testing.T) {
				testutils.ExpectGetUserByIdReturnsUser(t, queries)
				queries.EXPECT().DeleteOAuthAuthorization(t.Context(), models.DeleteOAuthAuthorizationParams{
					UserID:   testutils.MockUserId,
					Provider: models.ProviderGitHub,
				}).Return(nil).Once()
			},
		},
		{
			name:         "GitHub - with wrong password",
			verification: string(handlers.VerificationPassword),
			action:       handlers.ActionUnlinkGitHubAccount,
			password:     testutils.MockPassword + "wrong",
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnauthorized,
			before: func(t *testing.T) {
				testutils.ExpectGetUserByIdReturnsUser(t, queries)
			},
		},
		{
			name:         "GitHub - with google",
			verification: string(handlers.VerificationGoogle),
			action:       handlers.ActionUnlinkGitHubAccount,
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
			verification: string(handlers.VerificationGitHub),
			action:       handlers.ActionUnlinkGitHubAccount,
			password:     testutils.MockPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
			before: func(t *testing.T) {
				githubOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, "/user/oauth/github/callback/action").
					Return("https://github.com/oauth").Once()
			},
		},
		{
			name:         "Google - with password",
			verification: string(handlers.VerificationPassword),
			action:       handlers.ActionUnlinkGoogleAccount,
			password:     testutils.MockPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
			before: func(t *testing.T) {
				testutils.ExpectGetUserByIdReturnsUser(t, queries)
				queries.EXPECT().DeleteOAuthAuthorization(t.Context(), models.DeleteOAuthAuthorizationParams{
					UserID:   testutils.MockUserId,
					Provider: models.ProviderGoogle,
				}).Return(nil).Once()
			},
		},
		{
			name:         "Google - with wrong password",
			verification: string(handlers.VerificationPassword),
			action:       handlers.ActionUnlinkGoogleAccount,
			password:     testutils.MockPassword + "wrong",
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnauthorized,
			before: func(t *testing.T) {
				testutils.ExpectGetUserByIdReturnsUser(t, queries)
			},
		},
		{
			name:         "Google - with github",
			verification: string(handlers.VerificationGitHub),
			action:       handlers.ActionUnlinkGoogleAccount,
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
			verification: string(handlers.VerificationGoogle),
			action:       handlers.ActionUnlinkGoogleAccount,
			password:     testutils.MockPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
			before: func(t *testing.T) {
				googleOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, "/user/oauth/google/callback/action").
					Return("https://google.com/oauth").Once()
			},
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
			case handlers.ActionUnlinkGitHubAccount:
				form.Add("UnlinkGithub.Password", tt.password)
			case handlers.ActionUnlinkGoogleAccount:
				form.Add("UnlinkGoogle.Password", tt.password)
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

func TestUserUnlinkOAuthCallback(t *testing.T) {
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
		name             string
		verification     handlers.VerificationMethod
		wantCode         int
		wantRedirectPath string
		before           func(t *testing.T)
		after            func(t *testing.T)
	}{
		{
			name:             "Unlink GitHub account using Google verification",
			verification:     handlers.VerificationGoogle,
			wantCode:         http.StatusSeeOther,
			wantRedirectPath: "/user/settings",
			before: func(t *testing.T) {
				currentState = ts.MustSendUserSettingsOAuthAction(t, githubOauthService, googleOauthService, testutils.UserSettingsOAuthActionParams{
					Action:       handlers.ActionUnlinkGitHubAccount,
					Verification: handlers.VerificationGoogle,
					CSRFToken:    csrfToken,
				})

				queries.EXPECT().DeleteOAuthAuthorization(t.Context(), models.DeleteOAuthAuthorizationParams{
					UserID:   testutils.MockUserId,
					Provider: models.ProviderGitHub,
				}).Return(nil).Once()
			},
		},
		{
			name:             "Unlink GitHub account using GitHub verification",
			verification:     handlers.VerificationGitHub,
			wantCode:         http.StatusSeeOther,
			wantRedirectPath: "/user/settings",
			before: func(t *testing.T) {
				currentState = ts.MustSendUserSettingsOAuthAction(t, githubOauthService, googleOauthService, testutils.UserSettingsOAuthActionParams{
					Action:       handlers.ActionUnlinkGitHubAccount,
					Verification: handlers.VerificationGitHub,
					CSRFToken:    csrfToken,
				})

				queries.EXPECT().DeleteOAuthAuthorization(t.Context(), models.DeleteOAuthAuthorizationParams{
					UserID:   testutils.MockUserId,
					Provider: models.ProviderGitHub,
				}).Return(nil).Once()
			},
		},
		{
			name:             "Unlink Google account using GitHub verification",
			verification:     handlers.VerificationGitHub,
			wantCode:         http.StatusSeeOther,
			wantRedirectPath: "/user/settings",
			before: func(t *testing.T) {
				currentState = ts.MustSendUserSettingsOAuthAction(t, githubOauthService, googleOauthService, testutils.UserSettingsOAuthActionParams{
					Action:       handlers.ActionUnlinkGoogleAccount,
					Verification: handlers.VerificationGitHub,
					CSRFToken:    csrfToken,
				})

				queries.EXPECT().DeleteOAuthAuthorization(t.Context(), models.DeleteOAuthAuthorizationParams{
					UserID:   testutils.MockUserId,
					Provider: models.ProviderGoogle,
				}).Return(nil).Once()
			},
		},
		{
			name:             "Unlink Google account using Google verification",
			verification:     handlers.VerificationGoogle,
			wantCode:         http.StatusSeeOther,
			wantRedirectPath: "/user/settings",
			before: func(t *testing.T) {
				currentState = ts.MustSendUserSettingsOAuthAction(t, githubOauthService, googleOauthService, testutils.UserSettingsOAuthActionParams{
					Action:       handlers.ActionUnlinkGoogleAccount,
					Verification: handlers.VerificationGoogle,
					CSRFToken:    csrfToken,
				})

				queries.EXPECT().DeleteOAuthAuthorization(t.Context(), models.DeleteOAuthAuthorizationParams{
					UserID:   testutils.MockUserId,
					Provider: models.ProviderGoogle,
				}).Return(nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeEach(t)

			if tt.before != nil {
				tt.before(t)
			}

			res := ts.MustAuthenticateOAuthAction(t, testutils.AuthenticateOAuthActionParams{
				State:        currentState,
				Verification: tt.verification,
			})

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
