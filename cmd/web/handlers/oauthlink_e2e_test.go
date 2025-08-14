package handlers_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
	result := ts.Get(t, "/user/settings")
	assert.Equal(t, http.StatusOK, result.Status)
	csrfToken := testutils.ExtractCSRFToken(t, result.Body)
	assert.NotEmpty(t, csrfToken)

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
			action:       handlers.ActionLinkGitHubAccount,
			password:     password,
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
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
							String: testutils.MustHashPassword(t, password),
							Valid:  true,
						},
						Created: pgtype.Timestamptz{
							Time:  time.Now().Add(-time.Minute),
							Valid: true,
						},
					}, nil).Once()
				githubOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, "/user/oauth/github/callback/link").
					Return("https://github.com/oauth").Once()
			},
		},
		{
			name:         "GitHub - with wrong password",
			verification: string(handlers.VerificationPassword),
			action:       handlers.ActionLinkGitHubAccount,
			password:     password + "wrong",
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnauthorized,
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
							String: testutils.MustHashPassword(t, password),
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
			name:         "GitHub - with google",
			verification: string(handlers.VerificationGoogle),
			action:       handlers.ActionLinkGitHubAccount,
			password:     password,
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
			action:       handlers.ActionLinkGitHubAccount,
			password:     password,
			csrfToken:    csrfToken,
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Google - with password",
			verification: string(handlers.VerificationPassword),
			action:       handlers.ActionLinkGoogleAccount,
			password:     password,
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
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
							String: testutils.MustHashPassword(t, password),
							Valid:  true,
						},
						Created: pgtype.Timestamptz{
							Time:  time.Now().Add(-time.Minute),
							Valid: true,
						},
					}, nil).Once()
				googleOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, "/user/oauth/google/callback/link").
					Return("https://google.com/oauth").Once()
			},
		},
		{
			name:         "Google - with wrong password",
			verification: string(handlers.VerificationPassword),
			action:       handlers.ActionLinkGoogleAccount,
			password:     password + "wrong",
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnauthorized,
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
							String: testutils.MustHashPassword(t, password),
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
			name:         "Google - with github",
			verification: string(handlers.VerificationGitHub),
			action:       handlers.ActionLinkGoogleAccount,
			password:     password,
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
			action:       handlers.ActionLinkGoogleAccount,
			password:     password,
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
			case handlers.ActionLinkGitHubAccount:
				form.Add("LinkGithub.Password", tt.password)
			case handlers.ActionLinkGoogleAccount:
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

	var (
		username   = "walter"
		email      = "walter.white@example.com"
		password   = "LosPollos321"
		oauthCode  = "123456"
		oauthToken = "123456"
	)

	var currentState string

	ts.MustLogIn(t, queries, testutils.LoginUser{
		Username: username,
		Email:    email,
		Password: password,
	})

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
					Action:       handlers.ActionLinkGitHubAccount,
					Verification: handlers.VerificationGoogle,
					CSRFToken:    csrfToken,
				})

				ts.MustAuthenticateOAuthAction(t, testutils.AuthenticateOAuthActionParams{
					State:        currentState,
					Verification: handlers.VerificationGoogle,
					BeforeActionRedirect: func() {
						githubOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, mock.Anything).
							RunAndReturn(func(state string, callbackPath string) string {
								currentState = state
								return "https://github.com/oauth"
							}).Once()
					},
				})

				githubOauthService.EXPECT().ExchangeCodeForToken(services.TokenExchangeOptions{
					Code: oauthCode,
				}).Return(oauthToken, nil).Once()

				githubOauthService.EXPECT().GetUserInfo(oauthToken).Return(&services.OAuthUserInfo{
					Username: username,
					Id:       "1",
					Email:    email,
				}, nil).Once()

				queries.EXPECT().FindOAuthAuthorization(t.Context(), models.FindOAuthAuthorizationParams{
					UserProviderID: "1",
					Provider:       models.ProviderGitHub,
				}).Return(models.OauthAuthorization{}, pgx.ErrNoRows).Once()

				queries.EXPECT().CreateOAuthAuthorization(t.Context(), models.CreateOAuthAuthorizationParams{
					UserID:         1,
					Provider:       models.ProviderGitHub,
					UserProviderID: "1",
				}).Return(models.OauthAuthorization{
					ID:             1,
					UserID:         1,
					Provider:       models.ProviderGitHub,
					UserProviderID: "1",
				}, nil).Once()
			},
		},
		{
			name:             "Link Google account using GitHub verification",
			callbackPath:     "/user/oauth/google/callback/link",
			wantCode:         http.StatusSeeOther,
			wantRedirectPath: "/user/settings",
			before: func(t *testing.T) {
				currentState = ts.MustSendUserSettingsOAuthAction(t, githubOauthService, googleOauthService, testutils.UserSettingsOAuthActionParams{
					Action:       handlers.ActionLinkGoogleAccount,
					Verification: handlers.VerificationGitHub,
					CSRFToken:    csrfToken,
				})

				ts.MustAuthenticateOAuthAction(t, testutils.AuthenticateOAuthActionParams{
					State:        currentState,
					Verification: handlers.VerificationGitHub,
					BeforeActionRedirect: func() {
						googleOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, mock.Anything).
							RunAndReturn(func(state string, callbackPath string) string {
								currentState = state
								return "https://google.com/oauth"
							}).Once()
					},
				})

				googleOauthService.EXPECT().ExchangeCodeForToken(services.TokenExchangeOptions{
					Code:         oauthCode,
					RedirectPath: "/user/oauth/google/callback/link",
				}).Return(oauthToken, nil).Once()

				googleOauthService.EXPECT().GetUserInfo(oauthToken).Return(&services.OAuthUserInfo{
					Username: username,
					Id:       "1",
					Email:    email,
				}, nil).Once()

				queries.EXPECT().FindOAuthAuthorization(t.Context(), models.FindOAuthAuthorizationParams{
					UserProviderID: "1",
					Provider:       models.ProviderGoogle,
				}).Return(models.OauthAuthorization{}, pgx.ErrNoRows).Once()

				queries.EXPECT().CreateOAuthAuthorization(t.Context(), models.CreateOAuthAuthorizationParams{
					UserID:         1,
					Provider:       models.ProviderGoogle,
					UserProviderID: "1",
				}).Return(models.OauthAuthorization{
					ID:             1,
					UserID:         1,
					Provider:       models.ProviderGoogle,
					UserProviderID: "1",
				}, nil).Once()
			},
		},
		{
			name:             "GitHub account is already linked to an account",
			callbackPath:     "/user/oauth/github/callback/link",
			wantCode:         http.StatusSeeOther,
			wantRedirectPath: "/user/settings",
			before: func(t *testing.T) {
				currentState = ts.MustSendUserSettingsOAuthAction(t, githubOauthService, googleOauthService, testutils.UserSettingsOAuthActionParams{
					Action:       handlers.ActionLinkGitHubAccount,
					Verification: handlers.VerificationGoogle,
					CSRFToken:    csrfToken,
				})

				ts.MustAuthenticateOAuthAction(t, testutils.AuthenticateOAuthActionParams{
					State:        currentState,
					Verification: handlers.VerificationGoogle,
					BeforeActionRedirect: func() {
						githubOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, mock.Anything).
							RunAndReturn(func(state string, callbackPath string) string {
								currentState = state
								return "https://github.com/oauth"
							}).Once()
					},
				})

				githubOauthService.EXPECT().ExchangeCodeForToken(services.TokenExchangeOptions{
					Code: oauthCode,
				}).Return(oauthToken, nil).Once()

				githubOauthService.EXPECT().GetUserInfo(oauthToken).Return(&services.OAuthUserInfo{
					Username: username,
					Id:       "1",
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
			},
		},
		{
			name:             "Google account is already linked to an account",
			callbackPath:     "/user/oauth/google/callback/link",
			wantCode:         http.StatusSeeOther,
			wantRedirectPath: "/user/settings",
			before: func(t *testing.T) {
				currentState = ts.MustSendUserSettingsOAuthAction(t, githubOauthService, googleOauthService, testutils.UserSettingsOAuthActionParams{
					Action:       handlers.ActionLinkGoogleAccount,
					Verification: handlers.VerificationGitHub,
					CSRFToken:    csrfToken,
				})

				ts.MustAuthenticateOAuthAction(t, testutils.AuthenticateOAuthActionParams{
					State:        currentState,
					Verification: handlers.VerificationGitHub,
					BeforeActionRedirect: func() {
						googleOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, mock.Anything).
							RunAndReturn(func(state string, callbackPath string) string {
								currentState = state
								return "https://google.com/oauth"
							}).Once()
					},
				})

				googleOauthService.EXPECT().ExchangeCodeForToken(services.TokenExchangeOptions{
					Code:         oauthCode,
					RedirectPath: "/user/oauth/google/callback/link",
				}).Return(oauthToken, nil).Once()

				googleOauthService.EXPECT().GetUserInfo(oauthToken).Return(&services.OAuthUserInfo{
					Username: username,
					Id:       "1",
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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}

			res := ts.Get(t, tt.callbackPath+"?state="+currentState+"&code="+oauthCode)
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
