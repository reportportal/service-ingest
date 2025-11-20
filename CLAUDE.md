# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the ingest service for ReportPortal, responsible for receiving and processing agent data. It's a Go-based HTTP service built with the chi router framework.

## Development Commands

### Running the Service

- `make run` or `go run ./cmd/ingest/main.go` - Start the service locally on port 8080

### Building

- `go build ./...` - Compile all packages
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

- `cmd/ingest/main.go` - Main entry point, sets up chi router with middleware (Logger, Recoverer)
- `internal/` - Private application code
  - `data/` - Data layer (repositories, database interactions)
  - `domain/` - Domain models and business logic
- `pkg/` - Shared libraries that could be imported by external projects
- `api/` - API contracts (OpenAPI specs, DTOs, gRPC definitions)
- `bin/` - Compiled binaries

### HTTP Service

The service uses `go-chi/chi/v5` as the HTTP router. The main application:

- Runs on port 8080
- Has middleware for logging and panic recovery
- Currently has basic routes at `/` and `/panic` (likely placeholders)

## Coding Conventions

### Go Style

- Exported symbols: PascalCase
- Unexported symbols: camelCase
- Interfaces: Use `-er` suffix pattern (e.g., `Fooer`)
- Import grouping: stdlib, third-party, local (separated by blank lines)

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
