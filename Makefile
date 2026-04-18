.PHONY: dev build run deps

# Run server in development
dev:
	go run ./cmd/server/main.go

# Build binary
build:
	CGO_ENABLED=0 go build -o bin/server ./cmd/server/main.go

# Build and run
run: build
	./bin/server

# Install go dependencies
deps:
	go mod tidy
