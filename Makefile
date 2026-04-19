.PHONY: dev build run deps

# Run server in development
dev:
	go run ./cmd/api/main.go

# Build binary
build:
	CGO_ENABLED=0 go build -o bin/server ./cmd/api/main.go

# Build and run
run: build
	./bin/server

# Install go dependencies
deps:
	go mod tidy
