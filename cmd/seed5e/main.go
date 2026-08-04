package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/hadnu/onatar/internal/config"
	"github.com/hadnu/onatar/internal/etl5e"
	"github.com/hadnu/onatar/internal/etl5e/filter"
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

	fiveEToolsDir := os.Getenv("5ETOOLS_DIR")
	if fiveEToolsDir == "" {
		fiveEToolsDir = "vendor/5etools-src/data"
	}

	dryRun := os.Getenv("DRY_RUN") == "1"

	slog.Info("starting 5etools seed",
		"dir", fiveEToolsDir,
		"dry_run", dryRun,
	)

	// Load and normalize data
	data, err := etl5e.LoadFrom5eTools(fiveEToolsDir, filter.Filter2024Only)
	if err != nil {
		slog.Error("etl5e load", "error", err)
		os.Exit(1)
	}

	slog.Info("parsed 5etools data",
		"classes", len(data.Classes),
		"subclasses", len(data.Subclasses),
		"species", len(data.Species),
		"backgrounds", len(data.Backgrounds),
		"spells", len(data.Spells),
		"feats", len(data.Feats),
		"features", len(data.Features),
		"items", len(data.Items),
	)

	if dryRun {
		slog.Info("dry run complete, not writing to database")
		return
	}

	db, err := store.Open(cfg.DSN())
	if err != nil {
		slog.Error("connect db", "error", err)
		os.Exit(1)
	}
	defer func() { _ = db.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Full seed using existing store.Seed
	if err := store.Seed(ctx, db, data); err != nil {
		slog.Error("seed", "error", err)
		os.Exit(1)
	}

	slog.Info("seed5e complete",
		"classes", len(data.Classes),
		"subclasses", len(data.Subclasses),
		"species", len(data.Species),
		"backgrounds", len(data.Backgrounds),
		"spells", len(data.Spells),
		"feats", len(data.Feats),
		"features", len(data.Features),
		"items", len(data.Items),
	)
}