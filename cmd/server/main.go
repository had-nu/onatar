package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/hadnu/onatar/internal/config"
	"github.com/hadnu/onatar/internal/httpapi"
	"github.com/hadnu/onatar/internal/store"
)

const version = "0.1.0"

func main() {
	_ = godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "error", err)
		os.Exit(1)
	}

	db, err := store.Open(cfg.DSN())
	if err != nil {
		slog.Error("connect db", "error", err)
		os.Exit(1)
	}
	defer func() { _ = db.Close() }()

	srv := httpapi.New(db, logger, version)

	s := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           srv.Router(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	// #nosec G706 -- HTTPAddr comes from HTTP_ADDR env, operator-controlled.
	slog.Info("onatar server listening", "addr", cfg.HTTPAddr)
	if err := s.ListenAndServe(); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
