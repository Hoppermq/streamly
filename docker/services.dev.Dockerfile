FROM --platform=${BUILDPLATFORM} golang:1.25-alpine AS development
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app

# Install air for hot reload
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

# Use air for hot reload in development, fallback to go run if air config not found
ENTRYPOINT ["sh", "-c", "if [ -f \"$SERVICE_PATH/.air.toml\" ]; then cd $SERVICE_PATH && air; else go run $SERVICE_PATH; fi"]
