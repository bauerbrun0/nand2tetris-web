package main

import (
	"context"
	"database/sql"
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/application"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers/userhandlers"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/middleware"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/routes"
	"github.com/bauerbrun0/nand2tetris-web/db"
	"github.com/bauerbrun0/nand2tetris-web/internal"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/go-playground/form"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	err := godotenv.Load()
	if err != nil {
		logger.Info("Could not load .env file, proceesding with environment variables or flags", "error", err.Error())
	}
	var port int
	portEnv, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	} else {
		port = portEnv
	}

	var cfg application.Config
	flag.IntVar(&cfg.Port, "port", port, "HTTP server port")
	flag.BoolVar(&cfg.Migrate, "migrate", false, "Run database migrations")
	flag.StringVar(&cfg.Env, "env", os.Getenv("ENV"), "Environment (development|production )")
	flag.StringVar(&cfg.Dsn, "dsn", os.Getenv("DSN"), "Database Connection String")
	flag.StringVar(&cfg.BaseUrl, "base-url", os.Getenv("BASE_URL"), "The base URL of the application")
	flag.StringVar(&cfg.GithubClientId, "github-client-id", os.Getenv("GITHUB_CLIENT_ID"), "GitHub Client ID for OAuth")
	flag.StringVar(&cfg.GithubClientSecret, "github-client-secret", os.Getenv("GITHUB_CLIENT_SECRET"), "GitHub Client Secret for OAuth")
	flag.StringVar(&cfg.GoogleClientId, "google-client-id", os.Getenv("GOOGLE_CLIENT_ID"), "Google Client ID for OAuth")
	flag.StringVar(&cfg.GoogleClientSecret, "google-client-secret", os.Getenv("GOOGLE_CLIENT_SECRET"), "Google Client Secret for OAuth")
	flag.StringVar(&cfg.NoreplyEmail, "no-reply-email", os.Getenv("NOREPLY_EMAIL"), "No-reply email for sending emails to users")
	flag.StringVar(&cfg.MailgunDomain, "mailgun-domain", os.Getenv("MAILGUN_DOMAIN"), "Mailgun domain")
	flag.StringVar(&cfg.MailgunAPIKey, "mailgun-api-key", os.Getenv("MAILGUN_API_KEY"), "Mailgun private api key")
	flag.Parse()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.Dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer pool.Close()

	// Run database migrations if the migrate flag is set
	if cfg.Migrate {
		driver, err := postgres.WithInstance(sql.OpenDB(stdlib.GetPoolConnector(pool)), &postgres.Config{})
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		source, err := iofs.New(db.MigrationFiles, "migrations")
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		migrator, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		err = migrator.Up()

		if err != nil && errors.Is(err, migrate.ErrNoChange) {
			logger.Info("No database migrations to run")
			os.Exit(0)
		}

		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			logger.Error(err.Error())
			os.Exit(1)
		}

		logger.Info("Database migrations ran successfully")
		os.Exit(0)
	}

	gob.Register([]pages.Toast{})
	gob.Register(userhandlers.Action(""))
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(pool)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = cfg.Env == "production"

	queries := models.New(pool)
	txStarter := models.NewTxStarter(pool)

	var emailSender services.EmailSender

	if cfg.Env == "production" {
		emailSender = services.NewMailGunEmailSender(logger, cfg.MailgunDomain, cfg.MailgunAPIKey)
	} else {
		emailSender = services.NewConsoleEmailSender(logger)
	}

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	if cfg.Env == "production" {
		bundle.LoadMessageFileFS(internal.TranslationFiles, "translations/en.yaml")
	} else {
		bundle.LoadMessageFile("internal/translations/en.yaml")
	}

	emailService := services.NewEmailService(
		emailSender,
		logger,
		i18n.NewLocalizer(bundle, language.English.String()),
		cfg.NoreplyEmail,
		cfg.BaseUrl,
	)
	userService := services.NewUserService(logger, emailService, queries, txStarter, ctx)

	githubOauthService := services.NewGitHubOAuthService(cfg.GithubClientId, cfg.GithubClientSecret, cfg.BaseUrl, logger)
	googleOauthService := services.NewGoogleOAuthService(cfg.GoogleClientId, cfg.GoogleClientSecret, cfg.BaseUrl, logger)

	projectService := services.NewProjectService(logger, ctx, queries, txStarter)
	chipService := services.NewChipService(logger, ctx, queries, txStarter)

	app := &application.Application{
		Logger:             logger,
		Config:             cfg,
		SessionManager:     sessionManager,
		FormDecoder:        form.NewDecoder(),
		EmailService:       emailService,
		UserService:        userService,
		GithubOauthService: githubOauthService,
		GoogleOauthService: googleOauthService,
		ProjectService:     projectService,
		ChipService:        chipService,
		Bundle:             bundle,
	}

	handlers := handlers.NewHandlers(app)
	middleware := middleware.NewMiddleware(app)
	routes := routes.GetRoutes(app, middleware, handlers)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      routes,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Starting application", slog.String("env", cfg.Env), slog.Int("port", cfg.Port))

	err = srv.ListenAndServe()

	logger.Error(err.Error())
	os.Exit(1)
}
