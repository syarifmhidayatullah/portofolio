.PHONY: dev build run css docker-up docker-down migrate

# Run server in development
dev:
	go run ./cmd/server/main.go

# Build binary
build:
	go build -o bin/server ./cmd/server/main.go

# Build and run
run: build
	./bin/server

# Watch and compile Tailwind CSS
css:
	npx tailwindcss -i ./web/static/css/input.css -o ./web/static/css/app.css --watch

# Build Tailwind for production
css-build:
	npx tailwindcss -i ./web/static/css/input.css -o ./web/static/css/app.css --minify

# Start docker services
docker-up:
	docker-compose up -d

# Stop docker services
docker-down:
	docker-compose down

# Full docker build and up
docker-build:
	docker-compose up -d --build

# Install go dependencies
deps:
	go mod tidy
	npm install

# Run in dev (requires tmux or run in separate terminals)
# Terminal 1: make css
# Terminal 2: make dev
