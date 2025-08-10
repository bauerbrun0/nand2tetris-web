package handlers_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
	oauthCode := "12345678"
	oauthToken := "12345678"

	var (
		username = "walt"
		email    = "walter.white@example.com"
	)

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

				githubOauthService.EXPECT().ExchangeCodeForToken(services.TokenExchangeOptions{
					Code: oauthCode,
				}).Return(oauthToken, nil).Once()

				githubOauthService.EXPECT().GetUserInfo(oauthToken).Return(&services.OAuthUserInfo{
					Id:       "1",
					Username: username,
					Email:    email,
				}, nil).Once()

				queries.EXPECT().FindOAuthAuthorization(t.Context(), models.FindOAuthAuthorizationParams{
					UserProviderID: "1",
					Provider:       models.ProviderGitHub,
				}).Return(models.OauthAuthorization{
					ID:             1,
					UserID:         1,
					Provider:       models.ProviderGitHub,
					UserProviderID: "1",
				}, nil).Once()

				queries.EXPECT().GetUserInfo(t.Context(), int32(1)).Return(models.UserInfo{
					ID:       1,
					Username: username,
					Email:    email,
					EmailVerified: pgtype.Bool{
						Bool:  true,
						Valid: true,
					},
					Created: pgtype.Timestamptz{
						Time:  time.Now().Add(-time.Hour),
						Valid: true,
					},
					IsPasswordSet:  true,
					LinkedAccounts: []string{string(models.ProviderGitHub)},
				}, nil).Once()
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

				githubOauthService.EXPECT().ExchangeCodeForToken(services.TokenExchangeOptions{
					Code: oauthCode,
				}).Return(oauthToken, nil).Once()

				githubOauthService.EXPECT().GetUserInfo(oauthToken).Return(&services.OAuthUserInfo{
					Id:       "1",
					Username: username,
					Email:    email,
				}, nil).Once()

				queries.EXPECT().FindOAuthAuthorization(t.Context(), models.FindOAuthAuthorizationParams{
					UserProviderID: "1",
					Provider:       models.ProviderGitHub,
				}).Return(models.OauthAuthorization{}, pgx.ErrNoRows).Once()

				queries.EXPECT().GetUserInfoByEmailOrUsername(t.Context(), mock.Anything).
					Return(models.UserInfo{}, pgx.ErrNoRows).Once()

				queries.EXPECT().CreateNewUser(t.Context(), mock.Anything).Return(models.User{
					ID:       1,
					Username: username,
					Email:    email,
					EmailVerified: pgtype.Bool{
						Bool:  true,
						Valid: true,
					},
					PasswordHash: pgtype.Text{
						String: "string",
						Valid:  true,
					},
					Created: pgtype.Timestamptz{
						Time:  time.Now().Add(-time.Minute),
						Valid: true,
					},
				}, nil).Once()

				queries.EXPECT().CreateOAuthAuthorization(t.Context(), mock.Anything).Return(models.OauthAuthorization{
					ID:             1,
					UserID:         1,
					Provider:       models.ProviderGitHub,
					UserProviderID: "1",
				}, nil).Once()

				queries.EXPECT().GetUserInfo(t.Context(), int32(1)).Return(models.UserInfo{
					ID:       1,
					Username: username,
					Email:    email,
					EmailVerified: pgtype.Bool{
						Bool:  true,
						Valid: true,
					},
					Created: pgtype.Timestamptz{
						Time:  time.Now().Add(-time.Minute),
						Valid: true,
					},
					IsPasswordSet:  false,
					LinkedAccounts: []string{string(models.ProviderGitHub)},
				}, nil).Once()
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

				githubOauthService.EXPECT().ExchangeCodeForToken(services.TokenExchangeOptions{
					Code: oauthCode,
				}).Return(oauthToken, nil).Once()

				githubOauthService.EXPECT().GetUserInfo(oauthToken).Return(&services.OAuthUserInfo{
					Id:       "1",
					Username: username,
					Email:    email,
				}, nil).Once()

				queries.EXPECT().FindOAuthAuthorization(t.Context(), models.FindOAuthAuthorizationParams{
					UserProviderID: "1",
					Provider:       models.ProviderGitHub,
				}).Return(models.OauthAuthorization{}, pgx.ErrNoRows).Once()

				queries.EXPECT().GetUserInfoByEmailOrUsername(t.Context(), mock.Anything).
					Return(models.UserInfo{
						ID:       1,
						Username: username,
						Email:    email,
						EmailVerified: pgtype.Bool{
							Bool:  true,
							Valid: true,
						},
						Created: pgtype.Timestamptz{
							Time:  time.Now().Add(-time.Hour),
							Valid: true,
						},
						IsPasswordSet:  true,
						LinkedAccounts: []string{},
					}, nil).Once()
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

				googleOauthService.EXPECT().ExchangeCodeForToken(services.TokenExchangeOptions{
					Code:         oauthCode,
					RedirectPath: "/user/oauth/google/callback/login",
				}).Return(oauthToken, nil).Once()

				googleOauthService.EXPECT().GetUserInfo(oauthToken).Return(&services.OAuthUserInfo{
					Id:       "1",
					Username: username,
					Email:    email,
				}, nil).Once()

				queries.EXPECT().FindOAuthAuthorization(t.Context(), models.FindOAuthAuthorizationParams{
					UserProviderID: "1",
					Provider:       models.ProviderGoogle,
				}).Return(models.OauthAuthorization{
					ID:             1,
					UserID:         1,
					Provider:       models.ProviderGoogle,
					UserProviderID: "1",
				}, nil).Once()

				queries.EXPECT().GetUserInfo(t.Context(), int32(1)).Return(models.UserInfo{
					ID:       1,
					Username: username,
					Email:    email,
					EmailVerified: pgtype.Bool{
						Bool:  true,
						Valid: true,
					},
					Created: pgtype.Timestamptz{
						Time:  time.Now().Add(-time.Hour),
						Valid: true,
					},
					IsPasswordSet:  true,
					LinkedAccounts: []string{string(models.ProviderGoogle)},
				}, nil).Once()
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

				googleOauthService.EXPECT().ExchangeCodeForToken(services.TokenExchangeOptions{
					Code:         oauthCode,
					RedirectPath: "/user/oauth/google/callback/login",
				}).Return(oauthToken, nil).Once()

				googleOauthService.EXPECT().GetUserInfo(oauthToken).Return(&services.OAuthUserInfo{
					Id:       "1",
					Username: username,
					Email:    email,
				}, nil).Once()

				queries.EXPECT().FindOAuthAuthorization(t.Context(), models.FindOAuthAuthorizationParams{
					UserProviderID: "1",
					Provider:       models.ProviderGoogle,
				}).Return(models.OauthAuthorization{}, pgx.ErrNoRows).Once()

				queries.EXPECT().GetUserInfoByEmailOrUsername(t.Context(), mock.Anything).
					Return(models.UserInfo{}, pgx.ErrNoRows).Once()

				queries.EXPECT().CreateNewUser(t.Context(), mock.Anything).Return(models.User{
					ID:       1,
					Username: username,
					Email:    email,
					EmailVerified: pgtype.Bool{
						Bool:  true,
						Valid: true,
					},
					PasswordHash: pgtype.Text{
						String: "string",
						Valid:  true,
					},
					Created: pgtype.Timestamptz{
						Time:  time.Now().Add(-time.Minute),
						Valid: true,
					},
				}, nil).Once()

				queries.EXPECT().CreateOAuthAuthorization(t.Context(), mock.Anything).Return(models.OauthAuthorization{
					ID:             1,
					UserID:         1,
					Provider:       models.ProviderGoogle,
					UserProviderID: "1",
				}, nil).Once()

				queries.EXPECT().GetUserInfo(t.Context(), int32(1)).Return(models.UserInfo{
					ID:       1,
					Username: username,
					Email:    email,
					EmailVerified: pgtype.Bool{
						Bool:  true,
						Valid: true,
					},
					Created: pgtype.Timestamptz{
						Time:  time.Now().Add(-time.Minute),
						Valid: true,
					},
					IsPasswordSet:  false,
					LinkedAccounts: []string{string(models.ProviderGoogle)},
				}, nil).Once()
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

				googleOauthService.EXPECT().ExchangeCodeForToken(services.TokenExchangeOptions{
					Code:         oauthCode,
					RedirectPath: "/user/oauth/google/callback/login",
				}).Return(oauthToken, nil).Once()

				googleOauthService.EXPECT().GetUserInfo(oauthToken).Return(&services.OAuthUserInfo{
					Id:       "1",
					Username: username,
					Email:    email,
				}, nil).Once()

				queries.EXPECT().FindOAuthAuthorization(t.Context(), models.FindOAuthAuthorizationParams{
					UserProviderID: "1",
					Provider:       models.ProviderGoogle,
				}).Return(models.OauthAuthorization{}, pgx.ErrNoRows).Once()

				queries.EXPECT().GetUserInfoByEmailOrUsername(t.Context(), mock.Anything).
					Return(models.UserInfo{
						ID:       1,
						Username: username,
						Email:    email,
						EmailVerified: pgtype.Bool{
							Bool:  true,
							Valid: true,
						},
						Created: pgtype.Timestamptz{
							Time:  time.Now().Add(-time.Hour),
							Valid: true,
						},
						IsPasswordSet:  true,
						LinkedAccounts: []string{},
					}, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeEach(t)

			if tt.before != nil {
				tt.before(t)
			}

			path := fmt.Sprintf("%s?code=%s&state=%s", tt.callbackPath, oauthCode, currentState)
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
