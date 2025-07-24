package main

import (
	"context"
	"encoding/gob"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/bauerbrun0/nand2tetris-web/internal"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/go-playground/form"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type config struct {
	port int
	env  string
	dsn  string
}

type application struct {
	logger         *slog.Logger
	dev            bool
	config         config
	sessionManager *scs.SessionManager
	formDecoder    *form.Decoder
	emailService   *services.EmailService
	userService    *services.UserService
	bundle         *i18n.Bundle
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 3000, "HTTP server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|production )")
	flag.StringVar(&cfg.dsn, "dsn", "postgres://nand2tetris_web:password@localhost/nand2tetris_web?sslmode=disable", "Database Connection String")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

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

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	if cfg.env == "production" {
		bundle.LoadMessageFileFS(internal.TranslationFiles, "translations/en.yaml")
	} else {
		bundle.LoadMessageFile("internal/translations/en.yaml")
	}

	app := &application{
		logger:         logger,
		config:         cfg,
		sessionManager: sessionManager,
		formDecoder:    form.NewDecoder(),
		emailService:   emailService,
		userService:    userService,
		bundle:         bundle,
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
