# Repository Guidelines

## Project Structure & Module Organization

The ingest service centers on `cmd/ingest/main.go`, which wires up the CLI entry point. Place runtime-only binaries under `bin/`, and keep shared libraries in `pkg/`. Domain-specific packages should live in `internal/`, while API contracts (OpenAPI specs, DTOs, or gRPC stubs) belong in `api/`. Co-locate configuration samples or fixtures beside the code they exercise to keep modules discoverable.

## Build, Test, and Development Commands

- `go run ./cmd/ingest`: Run the ingest service locally; perfect for quick smoke checks.
- `go build -o bin/ ./...`: Compile all packages and surface type errors early.
- `go test ./...`: Execute the unit test suite; add `-race` when touching concurrency.
- `go vet ./...`: Catch suspicious patterns before review.
- `go fmt ./...`: Format code prior to committing; matches the repository baseline.

## Coding Style & Naming Conventions

Follow idiomatic Go style: exported types and funcs use PascalCase with meaningful names, private symbols stay camelCase, and interfaces adopt the `Fooer` pattern when appropriate. Let `gofmt` (tabs, K&R braces, trailing newlines) shape the file layout, and keep imports grouped standard/third-party/local. Prefer small, composable functions and return structured errors with context (`fmt.Errorf("context: %w", err)`).

## Testing Guidelines

Add `_test.go` files next to the code under test and rely on the standard `testing` package with table-driven cases. Cover happy paths, edge conditions, and failure scenarios whenever you modify behavior. Snapshot external responses under `testdata/` to avoid brittle network calls. Run `go test -race ./...` before review when touching goroutines or shared state.

## Commit & Pull Request Guidelines

Commits follow the concise, imperative style already in history (for example, “Add initial project files”). Keep titles under 72 characters, elaborate in the body when context matters, and group unrelated changes into separate commits. Pull requests should describe the change, call out impacted modules, link any ReportPortal issues, and include testing evidence (`go test`, manual verifications). Attach logs or screenshots when UI-visible outputs change so reviewers can validate quickly.

## Security & Configuration Tips

Never commit API tokens or production credentials; load them from environment variables or secret managers. Check new dependencies for licensing concerns and pin versions in `go.mod`. When adding outbound integrations, document required configuration keys in `README.md` and provide safe defaults for local development.
