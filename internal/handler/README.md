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
├── launch.go        # Launch-related endpoints
├── item.go          # Test item endpoints
├── log.go           # Log-related endpoints
├── info.go          # Service info endpoints
├── health.go        # Health check endpoints
├── router.go        # Route registration and middleware setup
├── middleware.go    # Custom HTTP middleware (auth, logging, etc.)
└── error.go         # Error handling models and utilities
```

## Router Setup

The `router.go` file sets up all routes and middleware.


## Handler models

Handlers use request and response models defined in the `*_dto.go` files to structure incoming and outgoing data.

## Dependencies

Handlers depend on:

- **service layer** - for business logic
- **model** - for domain models
- HTTP router (chi) - for routing and middleware

Handlers should NOT:

- Contain business logic (belongs in service layer)
- Access data layer directly (use service layer)
- Know about database details
