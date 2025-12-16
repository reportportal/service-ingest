package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/reportportal/service-ingest/internal/app"
	"github.com/reportportal/service-ingest/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	server, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to start new server: %v", err)
	}

	log.Printf("Starting server on %s (env: %s, base path: %s)",
		cfg.Server.Addr(),
		cfg.Server.Env,
		cfg.Server.BasePath,
	)

	if err := server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}
