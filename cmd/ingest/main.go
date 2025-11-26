package main

import (
	"log"
	"net/http"

	"github.com/reportportal/service-ingest/internal/config"
	"github.com/reportportal/service-ingest/internal/handler"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	router := handler.NewRouter(cfg.Server.BasePath)

	addr := cfg.Server.Addr()
	log.Printf("Starting server on %s (env: %s, base path: %s)", addr, cfg.Server.Env, cfg.Server.BasePath)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
