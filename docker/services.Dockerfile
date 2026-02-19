FROM --platform=${BUILDPLATFORM} golang:1.26-alpine AS deps
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download && go mod verify

FROM deps AS builder
ARG TARGET_OS
ARG SERVICE_PATH
ARG TARGET_ARCH
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 GOOS=${TARGET_OS} GOARCH=${TARGET_ARCH} \
  go build -ldflags='-w -s' -o service ${SERVICE_PATH}

FROM alpine:latest AS runtime
RUN apk --no-cache add ca-certificates tzdata && \
  adduser -D -s /bin/sh appuser
WORKDIR /app
COPY --from=builder /app/service ./service
USER appuser
ENTRYPOINT ["./service"]
