package userhandlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserLoginOAuth(t *testing.T) {
	ts, _, githubOauthService, googleOauthService := testutils.NewTestServer(t,
		testutils.TestServerOptions{
			Logs: false,
		},
	)
	defer ts.Close()

	var currentState string

	tests := []struct {
		name         string
		loginPath    string
		wantCode     int
		wantRedirect string
		before       func(t *testing.T)
		after        func(t *testing.T)
	}{
		{
			name:         "Valid GitHub",
			loginPath:    "/user/login/github",
			wantCode:     http.StatusSeeOther,
			wantRedirect: "https://github.com/oauth",
			before: func(t *testing.T) {
				githubOauthService.EXPECT().GetRedirectUrl(mock.Anything).
					RunAndReturn(func(state string) string {
						currentState = state
						return "https://github.com/oauth?state=" + state
					}).Once()
			},
		},
		{
			name:         "Valid Google",
			loginPath:    "/user/login/google",
			wantCode:     http.StatusSeeOther,
			wantRedirect: "https://google.com/oauth",
			before: func(t *testing.T) {
				googleOauthService.EXPECT().GetRedirectUrl(mock.Anything).
					RunAndReturn(func(state string) string {
						currentState = state
						return "https://google.com/oauth?state=" + state
					}).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}

			result := ts.Get(t, tt.loginPath)
			assert.Equal(t, tt.wantCode, result.Status)
			assert.Equal(t, tt.wantRedirect+"?state="+currentState, result.RedirectUrl.String())

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}

func TestUserLoginOAuthCallback(t *testing.T) {
	ts, queries, githubOauthService, googleOauthService := testutils.NewTestServer(t,
		testutils.TestServerOptions{
			Logs: false,
		},
	)
	defer ts.Close()

	var currentState string

	beforeEach := func(t *testing.T) {
		ts.RemoveCookie(t, "session")
	}

	tests := []struct {
		name         string
		callbackPath string
		wantCode     int
		wantRedirect string
		before       func(t *testing.T)
		after        func(t *testing.T)
	}{
		{
			name:         "GitHub - Existing user",
			callbackPath: "/user/oauth/github/callback/login",
			wantCode:     http.StatusSeeOther,
			wantRedirect: "/",
			before: func(t *testing.T) {
				githubOauthService.EXPECT().GetRedirectUrl(mock.Anything).
					RunAndReturn(func(state string) string {
						currentState = state
						return "https://github.com/oauth?state=" + state
					}).Once()
				ts.Get(t, "/user/login/github")

				testutils.ExpectExchangeCodeForUserInfo(t, githubOauthService)
				testutils.ExpectFindOAuthAuthorizationReturnsAuthorization(t, queries, models.ProviderGitHub)
				testutils.ExpectGetUserInfoReturnsUserInfoWithLinkedAccount(t, queries, models.ProviderGitHub)
			},
		},
		{
			name:         "GitHub - new user",
			callbackPath: "/user/oauth/github/callback/login",
			wantCode:     http.StatusSeeOther,
			wantRedirect: "/",
			before: func(t *testing.T) {
				githubOauthService.EXPECT().GetRedirectUrl(mock.Anything).
					RunAndReturn(func(state string) string {
						currentState = state
						return "https://github.com/oauth?state=" + state
					}).Once()
				ts.Get(t, "/user/login/github")

				testutils.ExpectExchangeCodeForUserInfo(t, githubOauthService)
				queries.EXPECT().FindOAuthAuthorization(t.Context(), models.FindOAuthAuthorizationParams{
					UserProviderID: testutils.MockOAuthUserId,
					Provider:       models.ProviderGitHub,
				}).Return(models.OauthAuthorization{}, pgx.ErrNoRows).Once()

				queries.EXPECT().GetUserInfoByEmailOrUsername(t.Context(), mock.Anything).
					Return(models.UserInfo{}, pgx.ErrNoRows).Once()

				testutils.ExpectCreateNewUserReturnsUser(t, queries)
				testutils.ExpectCreateOAuthAuthorizationReturnsAuthorization(t, queries, models.ProviderGitHub)
				testutils.ExpectGetUserInfoReturnsUserInfoWithLinkedAccount(t, queries, models.ProviderGitHub)
			},
		},
		{
			name:         "GitHub - email used already",
			callbackPath: "/user/oauth/github/callback/login",
			wantCode:     http.StatusSeeOther,
			wantRedirect: "/user/login",
			before: func(t *testing.T) {
				githubOauthService.EXPECT().GetRedirectUrl(mock.Anything).
					RunAndReturn(func(state string) string {
						currentState = state
						return "https://github.com/oauth?state=" + state
					}).Once()
				ts.Get(t, "/user/login/github")

				testutils.ExpectExchangeCodeForUserInfo(t, githubOauthService)

				queries.EXPECT().FindOAuthAuthorization(t.Context(), models.FindOAuthAuthorizationParams{
					UserProviderID: testutils.MockOAuthUserId,
					Provider:       models.ProviderGitHub,
				}).Return(models.OauthAuthorization{}, pgx.ErrNoRows).Once()
				testutils.ExpectGetUserInfoByEmailOrUsernameReturnsUserInfo(t, queries)
			},
		},
		{
			name:         "Google - Existing user",
			callbackPath: "/user/oauth/google/callback/login",
			wantCode:     http.StatusSeeOther,
			wantRedirect: "/",
			before: func(t *testing.T) {
				googleOauthService.EXPECT().GetRedirectUrl(mock.Anything).
					RunAndReturn(func(state string) string {
						currentState = state
						return "https://google.com/oauth?state=" + state
					}).Once()
				ts.Get(t, "/user/login/google")

				testutils.ExpectExchangeCodeForUserInfo(t, googleOauthService)

				testutils.ExpectFindOAuthAuthorizationReturnsAuthorization(t, queries, models.ProviderGoogle)
				testutils.ExpectGetUserInfoReturnsUserInfoWithLinkedAccount(t, queries, models.ProviderGoogle)
			},
		},
		{
			name:         "Google - new user",
			callbackPath: "/user/oauth/google/callback/login",
			wantCode:     http.StatusSeeOther,
			wantRedirect: "/",
			before: func(t *testing.T) {
				googleOauthService.EXPECT().GetRedirectUrl(mock.Anything).
					RunAndReturn(func(state string) string {
						currentState = state
						return "https://google.com/oauth?state=" + state
					}).Once()
				ts.Get(t, "/user/login/google")

				testutils.ExpectExchangeCodeForUserInfo(t, googleOauthService)

				queries.EXPECT().FindOAuthAuthorization(t.Context(), models.FindOAuthAuthorizationParams{
					UserProviderID: testutils.MockOAuthUserId,
					Provider:       models.ProviderGoogle,
				}).Return(models.OauthAuthorization{}, pgx.ErrNoRows).Once()

				queries.EXPECT().GetUserInfoByEmailOrUsername(t.Context(), mock.Anything).
					Return(models.UserInfo{}, pgx.ErrNoRows).Once()

				testutils.ExpectCreateNewUserReturnsUser(t, queries)
				testutils.ExpectCreateOAuthAuthorizationReturnsAuthorization(t, queries, models.ProviderGoogle)
				testutils.ExpectGetUserInfoReturnsUserInfoWithLinkedAccount(t, queries, models.ProviderGoogle)
			},
		},
		{
			name:         "Google - email used already",
			callbackPath: "/user/oauth/google/callback/login",
			wantCode:     http.StatusSeeOther,
			wantRedirect: "/user/login",
			before: func(t *testing.T) {
				googleOauthService.EXPECT().GetRedirectUrl(mock.Anything).
					RunAndReturn(func(state string) string {
						currentState = state
						return "https://google.com/oauth?state=" + state
					}).Once()
				ts.Get(t, "/user/login/google")

				testutils.ExpectExchangeCodeForUserInfo(t, googleOauthService)

				queries.EXPECT().FindOAuthAuthorization(t.Context(), models.FindOAuthAuthorizationParams{
					UserProviderID: testutils.MockOAuthUserId,
					Provider:       models.ProviderGoogle,
				}).Return(models.OauthAuthorization{}, pgx.ErrNoRows).Once()

				testutils.ExpectGetUserInfoByEmailOrUsernameReturnsUserInfo(t, queries)
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
			result := ts.Get(t, path)
			assert.Equal(t, tt.wantCode, result.Status)
			if tt.wantRedirect != "" {
				assert.Equal(t, tt.wantRedirect, result.RedirectUrl.Path)
			}

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}
