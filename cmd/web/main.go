package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
)

type application struct {
	logger         *slog.Logger
	dev            bool
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":3000", "HTTP network address")
	dev := flag.Bool("dev", false, "Development mode")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	sessionManager := scs.New()
	sessionManager.Store = memstore.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = *dev == false

	app := &application{
		logger:         logger,
		dev:            *dev,
		sessionManager: sessionManager,
	}

	logger.Info("Starting application", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, app.routes())

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
