package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/reportportal/service-ingest/internal/app"
	"github.com/reportportal/service-ingest/internal/config"
	"github.com/reportportal/service-ingest/pkg/logger"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err.Error()))
		os.Exit(1)
	}

	l := logger.New(logger.Options{
		Level:     cfg.Log.Level,
		Format:    cfg.Log.Format,
		AddSource: cfg.Log.AddSource,
	})
	slog.SetDefault(l)

	server, err := app.New(cfg)
	if err != nil {
		slog.Error("failed to start server", slog.Any("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("starting server",
		slog.String("service_name", "service-ingest"),
		slog.String("version", version),
		slog.String("commit", commit),
		slog.String("build_date", date),
		slog.String("base_path", cfg.Server.BasePath),
		slog.String("log_level", cfg.Log.Level),
		slog.String("address", cfg.Server.Address),
	)

	if err := server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server error", slog.Any("error", err.Error()))
		os.Exit(1)
	}
}
