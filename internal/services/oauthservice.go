package services

type OAuthUserInfo struct {
	Id        string
	Username  string
	Email     string
	AvatarUrl string
}

type OAuthService interface {
	GetRedirectUrl(state string) string
	ExchangeCodeForToken(code string) (string, error)
	GetUserInfo(token string) (*OAuthUserInfo, error)
}
