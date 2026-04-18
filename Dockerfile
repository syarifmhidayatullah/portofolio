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
COPY web/templates ./web/templates
COPY web/static ./web/static

EXPOSE 8080

CMD ["./server"]
