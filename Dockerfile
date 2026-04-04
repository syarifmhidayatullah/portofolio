FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build Tailwind (requires node)
FROM node:20-alpine AS css-builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY tailwind.config.js ./
COPY web/static/css/input.css ./web/static/css/
COPY web/templates ./web/templates
RUN npx tailwindcss -i ./web/static/css/input.css -o ./web/static/css/app.css --minify

# Build Go binary
FROM golang:1.23-alpine AS go-builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server/main.go

# Final image
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=go-builder /app/server .
COPY --from=css-builder /app/web/static/css/app.css ./web/static/css/app.css
COPY web/templates ./web/templates
COPY web/static ./web/static

EXPOSE 8080

CMD ["./server"]
