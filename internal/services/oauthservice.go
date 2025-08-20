package services

import "errors"

type OAuthUserInfo struct {
	Id        string
	Username  string
	Email     string
	AvatarUrl string
}

type TokenExchangeOptions struct {
	Code         string
	RedirectPath string
}

type OAuthService interface {
	GetRedirectUrl(state string) string
	ExchangeCodeForToken(options TokenExchangeOptions) (string, error)
	GetUserInfo(token string) (*OAuthUserInfo, error)
	GetRedirectUrlWithCustomCallbackPath(state string, callbackPath string) string
}

var ErrCouldNotGetOAuthUserInfo = errors.New("oauthservice: could not get oauth user info")
