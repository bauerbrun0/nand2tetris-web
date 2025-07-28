package services

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-resty/resty/v2"
)

type GoogleOAuthService struct {
	clientId     string
	clientSecret string
	logger       *slog.Logger
	client       *resty.Client
}

func NewGoogleOAuthService(clientId, clientSecret string, logger *slog.Logger) OAuthService {
	client := resty.New()
	return &GoogleOAuthService{
		clientId,
		clientSecret,
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
		"openid%20profile%20email",
		"http%3A//localhost%3A8080/user/login/google/callback",
	)
	return redirectUrl
}

type googleTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func (s *GoogleOAuthService) ExchangeCodeForToken(code string) (string, error) {
	data := map[string]string{
		"client_id":     s.clientId,
		"client_secret": s.clientSecret,
		"code":          code,
		"grant_type":    "authorization_code",
		"redirect_uri":  "http://localhost:8080/user/login/google/callback",
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
	_, err := s.client.R().
		SetAuthToken(token).
		SetHeader("Accept", "application/json").
		SetResult(googleUserResp).
		Get("https://www.googleapis.com/userinfo/v2/me")

	if err != nil {
		return nil, err
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
