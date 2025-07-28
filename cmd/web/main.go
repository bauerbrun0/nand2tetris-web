package main

import (
	"context"
	"encoding/gob"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/bauerbrun0/nand2tetris-web/internal"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/go-playground/form"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type config struct {
	port               int
	env                string
	dsn                string
	githubClientId     string
	githubClientSecret string
	googleClientId     string
	googleClientSecret string
}

type application struct {
	logger             *slog.Logger
	dev                bool
	config             config
	sessionManager     *scs.SessionManager
	formDecoder        *form.Decoder
	emailService       *services.EmailService
	userService        *services.UserService
	githubOauthService *services.GitHubOAuthService
	googleOauthService *services.GoogleOAuthService
	bundle             *i18n.Bundle
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	err := godotenv.Load()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	var cfg config
	flag.IntVar(&cfg.port, "port", port, "HTTP server port")
	flag.StringVar(&cfg.env, "env", os.Getenv("ENV"), "Environment (development|production )")
	flag.StringVar(&cfg.dsn, "dsn", os.Getenv("DSN"), "Database Connection String")
	flag.StringVar(&cfg.githubClientId, "github-client-id", os.Getenv("GITHUB_CLIENT_ID"), "GitHub Client ID for OAuth")
	flag.StringVar(&cfg.githubClientSecret, "github-client-secret", os.Getenv("GITHUB_CLIENT_SECRET"), "GitHub Client Secret for OAuth")
	flag.StringVar(&cfg.googleClientId, "google-client-id", os.Getenv("GOOGLE_CLIENT_ID"), "Google Client ID for OAuth")
	flag.StringVar(&cfg.googleClientSecret, "google-client-secret", os.Getenv("GOOGLE_CLIENT_SECRET"), "Google Client Secret for OAuth")
	flag.Parse()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer pool.Close()

	gob.Register([]pages.Toast{})
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(pool)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = cfg.env == "production"

	emailSender := services.NewConsoleEmailSender(logger)
	emailService := services.NewEmailService(emailSender, logger)
	userService := services.NewUserService(logger, emailService, pool, ctx)

	githubOauthService := services.NewGitHubOAuthService(cfg.githubClientId, cfg.githubClientSecret, logger)
	googleOauthService := services.NewGoogleOAuthService(cfg.googleClientId, cfg.googleClientSecret, logger)

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	if cfg.env == "production" {
		bundle.LoadMessageFileFS(internal.TranslationFiles, "translations/en.yaml")
	} else {
		bundle.LoadMessageFile("internal/translations/en.yaml")
	}

	app := &application{
		logger:             logger,
		config:             cfg,
		sessionManager:     sessionManager,
		formDecoder:        form.NewDecoder(),
		emailService:       emailService,
		userService:        userService,
		githubOauthService: githubOauthService,
		googleOauthService: googleOauthService,
		bundle:             bundle,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Starting application", slog.String("env", cfg.env), slog.Int("port", cfg.port))

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
