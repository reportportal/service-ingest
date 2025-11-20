.PHONY: run build clean

# Run the service directly with go run (for development)
run:
	go run ./cmd/ingest/main.go

# Build the binary
build:
	go build -o bin/ingest ./cmd/ingest

# Clean build artifacts
clean:
	rm -rf bin/