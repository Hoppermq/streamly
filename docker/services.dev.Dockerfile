FROM --platform=${BUILDPLATFORM} golang:1.24-alpine AS development
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

ENTRYPOINT ["sh", "-c", "go run $SERVICE_PATH"]
