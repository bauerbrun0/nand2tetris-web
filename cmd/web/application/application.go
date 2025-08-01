package application

import (
	"log/slog"

	"github.com/alexedwards/scs/v2"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/go-playground/form"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Config struct {
	Port               int
	Env                string
	Dsn                string
	BaseUrl            string
	GithubClientId     string
	GithubClientSecret string
	GoogleClientId     string
	GoogleClientSecret string
}

type Application struct {
	Logger             *slog.Logger
	Dev                bool
	Config             Config
	SessionManager     *scs.SessionManager
	FormDecoder        *form.Decoder
	EmailService       *services.EmailService
	UserService        *services.UserService
	GithubOauthService services.OAuthService
	GoogleOauthService services.OAuthService
	Bundle             *i18n.Bundle
}
