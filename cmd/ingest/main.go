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

	router := handler.NewRouter()

	addr := cfg.Server.Addr()
	log.Printf("Starting server on %s (env: %s)", addr, cfg.Server.Env)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
