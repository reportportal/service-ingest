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
├── launch.go        # Launch-related endpoints and LaunchHandler type
├── launch_model.go  # Launch request/response models
├── item.go          # Test item endpoints and ItemHandler type
├── item_model.go    # Item request/response models
├── log.go           # Log-related endpoints and LogHandler type
├── log_model.go     # Log request/response models
├── info.go          # Service info endpoints
├── health.go        # Health check endpoints
├── route.go         # Router setup, Handlers struct, and route registration
├── middleware.go    # Custom HTTP middleware (parsing, validation, etc.)
├── validator.go     # Validator singleton and configuration
└── error.go         # Error handling models and utilities
```

## Handler Initialization

Each handler is a struct with dependencies injected via constructor:

```go
type LaunchHandler struct {
    service *service.LaunchService
}

func NewLaunchHandler(svc *service.LaunchService) *LaunchHandler {
    return &LaunchHandler{service: svc}
}
```

Similarly, for `ItemHandler` and `LogHandler`.

## Handlers Container

All handlers are grouped in a `Handlers` struct for easier dependency management:

```go
type Handlers struct {
    Launch *LaunchHandler
    Item   *ItemHandler
    Log    *LogHandler
}
```

This container is passed to `NewRouter()` during initialization.

## Router Setup

The `route.go` file contains:
- `Handlers` struct - container for all handlers
- `NewRouter(basePath string, handlers Handlers)` - sets up routes and middleware
- `apiRouter(handlers Handlers)` - configures API routes

Example usage from `internal/app/app.go`:

```go
handlers := handler.Handlers{
    Launch: handler.NewLaunchHandler(launchService),
    Item:   handler.NewItemHandler(itemService),
    Log:    handler.NewLogHandler(logService),
}

router := handler.NewRouter(cfg.Server.BasePath, handlers)
```

## Handler Models

Handlers use request and response models defined in `*_model.go` files to structure incoming and outgoing data.

## Validation

The `validator.go` file provides a singleton validator instance using `github.com/go-playground/validator/v10`:

```go
var validate = validator.New()
```

The validator is configured to use JSON field names in validation error messages instead of struct field names. This is initialized once using `sync.Once` pattern via `initValidatorOnce()` which is called during router setup.

Model structs use validation tags:

```go
type StartLaunchRQ struct {
    Name string `json:"name" validate:"required"`
    UUID string `json:"uuid" validate:"omitempty,uuid"`
}
```

Validation is performed in request binding methods using the shared validator instance.

## Dependencies

Handlers depend on:

- **service layer** - for business logic
- **model** - for domain models
- HTTP router (chi) - for routing and middleware

Handlers should NOT:

- Contain business logic (belongs in service layer)
- Access data layer directly (use service layer)
- Know about database details
