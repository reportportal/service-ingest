# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the ingest service for ReportPortal, responsible for receiving and processing agent data. It's a Go-based HTTP service built with the chi router framework.

## Development Commands

### Running the Service

- `make run` or `go run ./cmd/ingest/main.go` - Start the service locally on port 8080

### Building

- `go build -o bin/ ./...` - Compile all packages
- `go build -o bin/ingest ./cmd/ingest` - Build the main binary

### Running tests

- `go test ./...` - Run all tests
- `go test -race ./...` - Run tests with race detection (use when modifying concurrent code)
- `go test ./path/to/package` - Run tests for a specific package

### Code Quality

- `go vet ./...` - Static analysis for suspicious code
- `go fmt ./...` - Format code (run before committing)

## Architecture

### Project Structure

The service follows a three-layer architecture: **handler → service → data**

- `cmd/ingest/main.go` - Entry point, dependency injection, wires up all layers
- `internal/` - Private application code organized in layers:
  - `handler/` - HTTP handlers (controllers), request/response handling
  - `service/` - Business logic layer, orchestrates operations
  - `data/` - Data access layer (repositories, database interactions)
  - `model/` - Domain models (entities) with validation
  - `config/` - Configuration management from environment variables
- `pkg/` - Shared libraries that could be imported by external projects
- `api/` - API contracts (OpenAPI specs, DTOs, gRPC definitions)
- `migrations/` - Database migration scripts
- `bin/` - Compiled binaries

Each layer has its own README.md with detailed documentation and examples.

### Layer Responsibilities

**handler/** - HTTP layer (presentation)

- Parse and validate HTTP requests
- Call service layer methods
- Format and return HTTP responses
- Handle HTTP-specific concerns (status codes, headers)

**service/** - Business logic layer

- Implement business rules and validation
- Coordinate between multiple repositories
- Handle transactions and orchestration
- Return domain errors

**data/** - Data access layer

- Implement CRUD operations
- Execute database queries
- Manage connections (PostgreSQL, Redis, etc.)
- No business logic

**model/** - Domain models

- Define entities with validation
- Pure domain logic, no dependencies
- Used across all layers

**config/** - Configuration

- Load from environment variables
- Provide defaults and validation
- Simple approach without external libraries

### HTTP Service

The service uses `go-chi/chi/v5` as the HTTP router:

- Runs on port 8080 (configurable via `PORT` env variable)
- Middleware: Logger, Recoverer
- Routes defined in `internal/handler/router.go`
- API follows pattern: `/api/v1/{resource}`

## Coding Conventions

### Go Style

- Exported symbols: PascalCase
- Unexported symbols: camelCase
- Interfaces: Use `-er` suffix pattern (e.g., `Fooer`) and add them when needed, not upfront
- Import grouping: stdlib, third-party, local (separated by blank lines)

### Architecture Principles

- **Start simple** - Use concrete types; add interfaces only when mocking or abstraction is needed
- **Explicit dependencies** - Constructor injection, no magic DI containers
- **Layer isolation** - Handlers don't access data directly; services orchestrate
- **Go idioms first** - Favor Go conventions over framework patterns
- **Documentation** - Each layer has README.md with examples

### Error Handling

Return structured errors with context using `fmt.Errorf("context: %w", err)` to maintain error chains.

### Testing

- Place `_test.go` files alongside the code being tested
- Use table-driven tests
- Store test fixtures and mock responses in `testdata/` directories
- Use `-race` flag when testing concurrent code

### Commits

- Use imperative mood ("Add feature" not "Added feature")
- Keep commit titles under 72 characters
- Reference ReportPortal issues when applicable

## Configuration

Environment variables and secrets should never be committed. Use `.env.example` as a template for local development configuration.
