.PHONY: dev build run css css-build deps

# Run server in development
dev:
	go run ./cmd/server/main.go

# Build binary
build:
	CGO_ENABLED=0 go build -o bin/server ./cmd/server/main.go

# Build and run
run: build
	./bin/server

# Watch and compile Tailwind CSS
css:
	npx tailwindcss -i ./web/static/css/input.css -o ./web/static/css/app.css --watch

# Build Tailwind for production
css-build:
	npx tailwindcss -i ./web/static/css/input.css -o ./web/static/css/app.css --minify

# Install dependencies
deps:
	go mod tidy
	npm install
