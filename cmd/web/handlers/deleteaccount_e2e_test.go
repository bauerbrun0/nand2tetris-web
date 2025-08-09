package handlers_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers"
	"github.com/bauerbrun0/nand2tetris-web/internal/crypto"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleUserSettingsDeleteAccountPost(t *testing.T) {
	ts, queries, githubOauthService, googleOauthService := testutils.NewTestServer(
		t,
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
		wantCode     int
		before       func(t *testing.T)
		after        func(t *testing.T)
	}{
		{
			name:         "With password verification",
			email:        email,
			verification: string(handlers.VerificationPassword),
			password:     password,
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
			before: func(t *testing.T) {
				var hasher crypto.PasswordHasher
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
							String: testutils.MustHashPassword(t, hasher, password),
							Valid:  true,
						},
						Created: pgtype.Timestamptz{
							Time:  time.Now().Add(-time.Minute),
							Valid: true,
						},
					}, nil).Once()
				queries.EXPECT().DeleteUser(t.Context(), int32(1)).
					Return(nil).Once()
			},
			after: func(t *testing.T) {
				ts.MustLogIn(t, queries, testutils.LoginUser{
					Username: username,
					Email:    email,
					Password: password,
				})
			},
		},
		{
			name:         "With GitHub verification",
			email:        email,
			verification: string(handlers.VerificationGitHub),
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
			before: func(t *testing.T) {
				githubOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, mock.Anything).
					Return(ts.URL).Once()
			},
			after: func(t *testing.T) {
				ts.MustLogIn(t, queries, testutils.LoginUser{
					Username: username,
					Email:    email,
					Password: password,
				})
			},
		},
		{
			name:         "With Google verification",
			email:        email,
			verification: string(handlers.VerificationGoogle),
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
			before: func(t *testing.T) {
				googleOauthService.EXPECT().GetRedirectUrlWithCustomCallbackPath(mock.Anything, mock.Anything).
					Return(ts.URL).Once()
			},
			after: func(t *testing.T) {
				ts.MustLogIn(t, queries, testutils.LoginUser{
					Username: username,
					Email:    email,
					Password: password,
				})
			},
		},
		{
			name:         "Empty email",
			email:        "",
			verification: string(handlers.VerificationPassword),
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
		},
		{
			name:         "Wrong email",
			email:        "wrong" + email,
			verification: string(handlers.VerificationPassword),
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
		},
		{
			name:         "Invalid verification",
			email:        email,
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
			form.Add("Action", string(handlers.ActionDeleteAccount))
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
	ts, queries, githubOauthService, googleOauthService := testutils.NewTestServer(
		t,
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

	beforeEach := func(t *testing.T) {
		ts.MustLogIn(t, queries, testutils.LoginUser{
			Username: username,
			Email:    email,
			Password: password,
		})
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
					Action:       handlers.ActionDeleteAccount,
					Verification: handlers.VerificationGitHub,
					CSRFToken:    csrfToken,
					FormData: map[string]string{
						"DeleteAccount.Email": email,
					},
				})

				githubOauthService.EXPECT().ExchangeCodeForToken(services.TokenExchangeOptions{
					Code: oauthCode,
				}).
					Return(oauthToken, nil).Once()

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

				queries.EXPECT().DeleteUser(t.Context(), int32(1)).Return(nil).Once()
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
					Action:       handlers.ActionDeleteAccount,
					Verification: handlers.VerificationGoogle,
					CSRFToken:    csrfToken,
					FormData: map[string]string{
						"DeleteAccount.Email": email,
					},
				})

				googleOauthService.EXPECT().ExchangeCodeForToken(services.TokenExchangeOptions{
					Code:         oauthCode,
					RedirectPath: "/user/oauth/google/callback/action",
				}).
					Return(oauthToken, nil).Once()

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

				queries.EXPECT().DeleteUser(t.Context(), int32(1)).Return(nil).Once()
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
					Action:       handlers.ActionDeleteAccount,
					Verification: handlers.VerificationGitHub,
					CSRFToken:    csrfToken,
					FormData: map[string]string{
						"DeleteAccount.Email": email,
					},
				})

				githubOauthService.EXPECT().ExchangeCodeForToken(services.TokenExchangeOptions{
					Code: oauthCode,
				}).
					Return(oauthToken, nil).Once()

				githubOauthService.EXPECT().GetUserInfo(oauthToken).Return(&services.OAuthUserInfo{
					Id:       "1",
					Username: username,
					Email:    email,
				}, nil).Once()

				queries.EXPECT().FindOAuthAuthorization(t.Context(), models.FindOAuthAuthorizationParams{
					UserProviderID: "1",
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

			path := fmt.Sprintf("%s?code=%s&state=%s", tt.callbackPath, oauthCode, currentState)
			result = ts.Get(t, path)
			assert.Equal(t, tt.wantCode, result.Status)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}
