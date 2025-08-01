package services

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
)

type GoogleOAuthService struct {
	clientId     string
	clientSecret string
	appBaseUrl   string
	logger       *slog.Logger
	client       *resty.Client
}

func NewGoogleOAuthService(clientId, clientSecret, appBaseUrl string, logger *slog.Logger) OAuthService {
	client := resty.New()
	return &GoogleOAuthService{
		clientId,
		clientSecret,
		appBaseUrl,
		logger,
		client,
	}
}

func (s *GoogleOAuthService) GetRedirectUrl(state string) string {
	redirectUrl := fmt.Sprintf(
		"%s?response_type=code&client_id=%s&state=%s&scope=%s&redirect_uri=%s",
		"https://accounts.google.com/o/oauth2/v2/auth",
		s.clientId,
		state,
		url.QueryEscape("openid profile email"),
		url.QueryEscape(fmt.Sprintf("%s/user/oauth/google/callback/login", s.appBaseUrl)),
	)
	return redirectUrl
}

func (s *GoogleOAuthService) GetRedirectUrlWithCustomCallbackPath(state string, callbackPath string) string {
	redirectUrl := fmt.Sprintf(
		"%s?response_type=code&client_id=%s&state=%s&scope=%s&redirect_uri=%s",
		"https://accounts.google.com/o/oauth2/v2/auth",
		s.clientId,
		state,
		url.QueryEscape("openid profile email"),
		url.QueryEscape(fmt.Sprintf("%s%s", s.appBaseUrl, callbackPath)),
	)
	return redirectUrl
}

type googleTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func (s *GoogleOAuthService) ExchangeCodeForToken(options TokenExchangeOptions) (string, error) {
	data := map[string]string{
		"client_id":     s.clientId,
		"client_secret": s.clientSecret,
		"code":          options.Code,
		"grant_type":    "authorization_code",
		"redirect_uri":  fmt.Sprintf("%s%s", s.appBaseUrl, options.RedirectPath),
	}

	googleResp := &googleTokenResponse{}

	_, err := s.client.R().
		SetFormData(data).
		SetHeader("Accept", "application/json").
		SetResult(googleResp).
		Post("https://oauth2.googleapis.com/token")

	if err != nil {
		return "", err
	}

	return googleResp.AccessToken, nil
}

type googleUserInfoResponse struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

func (s *GoogleOAuthService) GetUserInfo(token string) (*OAuthUserInfo, error) {
	googleUserResp := &googleUserInfoResponse{}
	resp, err := s.client.R().
		SetAuthToken(token).
		SetHeader("Accept", "application/json").
		SetResult(googleUserResp).
		Get("https://www.googleapis.com/userinfo/v2/me")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, ErrCouldNotGetOAuthUserInfo
	}

	username := strings.Split(googleUserResp.Email, "@")[0]

	userInfo := &OAuthUserInfo{
		Id:        googleUserResp.Id,
		Username:  username,
		Email:     googleUserResp.Email,
		AvatarUrl: googleUserResp.Picture,
	}
	return userInfo, nil
}
