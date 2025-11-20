# Handlers

This module contains the HTTP handlers (controllers) for the service-ingest application.
Handlers are responsible for processing HTTP requests, validating input, calling business logic,
and returning HTTP responses.

## Responsibilities

- Parse and validate HTTP requests
- Extract parameters from URL, query strings, and request body
- Call appropriate service layer methods
- Handle errors and return appropriate HTTP status codes
- Format and return HTTP responses (JSON, etc.)

## Structure

```text
handler/
├── agent.go         # Agent-related endpoints
├── launch.go        # Launch-related endpoints
├── test.go          # Test item endpoints
├── router.go        # Route registration and middleware setup
├── middleware.go    # Custom HTTP middleware (auth, logging, etc.)
└── response.go      # Helper functions for HTTP responses
```

## Example Handler

```go
package handler

import (
    "encoding/json"
    "net/http"
    "github.com/reportportal/service-ingest/internal/model"
    "github.com/reportportal/service-ingest/internal/service"
)

type AgentHandler struct {
    agentService *service.AgentService
}

func NewAgentHandler(svc *service.AgentService) *AgentHandler {
    return &AgentHandler{agentService: svc}
}

func (h *AgentHandler) Register(w http.ResponseWriter, r *http.Request) {
    var agent model.Agent
    if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
        respondError(w, http.StatusBadRequest, "invalid request")
        return
    }

    if err := h.agentService.Register(r.Context(), &agent); err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondJSON(w, http.StatusCreated, agent)
}
```

## Router Setup

The `router.go` file sets up all routes and middleware:

```go
package handler

import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(agentHandler *AgentHandler) *chi.Mux {
    r := chi.NewRouter()

    // Middleware
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    // Routes
    r.Route("/api/v1", func(r chi.Router) {
        r.Route("/agents", func(r chi.Router) {
            r.Post("/", agentHandler.Register)
            r.Get("/{id}", agentHandler.GetByID)
        })
    })

    return r
}
```

## Dependencies

Handlers depend on:

- **service layer** - for business logic
- **model** - for domain models
- HTTP router (chi) - for routing and middleware

Handlers should NOT:

- Contain business logic (belongs in service layer)
- Access data layer directly (use service layer)
- Know about database details
