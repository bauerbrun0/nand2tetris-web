package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
	dev    bool
}

func main() {
	addr := flag.String("addr", ":3000", "HTTP network address")
	dev := flag.Bool("dev", false, "Development mode")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		logger: logger,
		dev:    *dev,
	}

	logger.Info("Starting application", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, app.routes())

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
