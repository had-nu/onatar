package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/hadnu/onatar/internal/config"
	"github.com/hadnu/onatar/internal/content"
	"github.com/hadnu/onatar/internal/store"
)

func main() {
	_ = godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "error", err)
		os.Exit(1)
	}

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}

	data, err := content.LoadData(dataDir)
	if err != nil {
		slog.Error("load content", "error", err)
		os.Exit(1)
	}

	db, err := store.Open(cfg.DSN())
	if err != nil {
		slog.Error("connect db", "error", err)
		os.Exit(1)
	}
	defer func() { _ = db.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := store.Seed(ctx, db, data); err != nil {
		slog.Error("seed", "error", err)
		os.Exit(1)
	}

	// #nosec G706 -- counts are derived from local repo data, not user input.
	slog.Info("seed complete",
		"classes", len(data.Classes),
		"subclasses", len(data.Subclasses),
		"species", len(data.Species),
		"backgrounds", len(data.Backgrounds),
		"spells", len(data.Spells),
		"feats", len(data.Feats),
		"features", len(data.Features),
	)
}
