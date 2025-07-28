package services

import (
	"fmt"
	"log/slog"
	"net/url"
	"strconv"

	"github.com/go-resty/resty/v2"
)

type GitHubOAuthService struct {
	clientId     string
	clientSecret string
	appBaseUrl   string
	logger       *slog.Logger
	client       *resty.Client
}

func NewGitHubOAuthService(clientId, clientSecret, appBaseUrl string, logger *slog.Logger) OAuthService {
	client := resty.New()
	return &GitHubOAuthService{
		clientId,
		clientSecret,
		appBaseUrl,
		logger,
		client,
	}
}

func (s *GitHubOAuthService) GetRedirectUrl(state string) string {
	redirectUrl := fmt.Sprintf(
		"%s?client_id=%s&state=%s&scope=%s&redirect_uri=%s",
		"https://github.com/login/oauth/authorize",
		s.clientId,
		state,
		url.QueryEscape("user:email"),
		url.QueryEscape(fmt.Sprintf("%s/user/login/github/callback", s.appBaseUrl)),
	)
	return redirectUrl
}

type githubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func (s *GitHubOAuthService) ExchangeCodeForToken(code string) (string, error) {
	requestBody := map[string]string{
		"client_id":     s.clientId,
		"client_secret": s.clientSecret,
		"code":          code,
	}

	ghresp := &githubTokenResponse{}

	_, err := s.client.R().
		SetBody(requestBody).
		SetHeader("Accept", "application/json").
		SetResult(ghresp).
		Post("https://github.com/login/oauth/access_token")

	if err != nil {
		return "", err
	}

	return ghresp.AccessToken, nil
}

type githubUserInfoResponse struct {
	Id        int    `json:"id"`
	Username  string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
}

type githubUserEmailsResponse []struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

func (s *GitHubOAuthService) GetUserInfo(token string) (*OAuthUserInfo, error) {
	ghUserResp := &githubUserInfoResponse{}
	_, err := s.client.R().
		SetAuthToken(token).
		SetHeader("Accept", "application/json").
		SetResult(ghUserResp).
		Get("https://api.github.com/user")

	if err != nil {
		return nil, err
	}

	ghUserEmailsResp := &githubUserEmailsResponse{}
	_, err = s.client.R().
		SetAuthToken(token).
		SetHeader("Accept", "application/json").
		SetResult(ghUserEmailsResp).
		Get("https://api.github.com/user/emails")

	if err != nil {
		return nil, err
	}

	var email string
	for _, e := range *ghUserEmailsResp {
		if e.Primary {
			email = e.Email
			break
		}
	}

	userInfo := &OAuthUserInfo{
		Id:        strconv.Itoa(ghUserResp.Id),
		Username:  ghUserResp.Username,
		AvatarUrl: ghUserResp.AvatarUrl,
		Email:     email,
	}
	return userInfo, nil
}
