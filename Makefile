VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

LDFLAGS := -ldflags "\
	-X main.version=$(VERSION) \
  -X main.commit=$(COMMIT) \
  -X main.buildTime=$(BUILD_TIME)"

.PHONY: run build clean

# Run the service directly with go run (for development)
run:
	@echo "Building ingest service..."
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	go run $(LDFLAGS) ./cmd/ingest/main.go

# Build the binary
build:
	echo "Building ingest service..."
	echo "Version: $(VERSION)"
	echo "Commit: $(COMMIT)"
	echo "Build Time: $(BUILD_TIME)"
	go build $(LDFLAGS) -o bin/ingest ./cmd/ingest

# Clean build artifacts
clean:
	rm -rf bin/