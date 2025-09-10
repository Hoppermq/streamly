ARG BUILD_TARGET=production
ARG BUN_VERSION=1.0

FROM oven/bun:${BUN_VERSION}-alpine AS base
WORKDIR /app

RUN apk add --no-cache git

COPY frontend/package.json frontend/bun.lockb* ./

FROM base AS development
ENV NODE_ENV=development

RUN --mount=type=cache,target=/root/.bun \
    bun install --frozen-lockfile

COPY frontend/ .

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:3000/ || exit 1

CMD ["bun", "run", "dev"]

FROM base AS builder
ENV NODE_ENV=production

RUN --mount=type=cache,target=/root/.bun \
    bun install --frozen-lockfile

COPY frontend/ .

RUN bun run build

RUN ls -la dist/

FROM base AS production
ENV NODE_ENV=production

RUN --mount=type=cache,target=/root/.bun \
    bun install --frozen-lockfile --production

COPY --from=builder /app/dist ./dist

EXPOSE 3000
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:3000/ || exit 1
CMD ["bun", "run", "serve"]
FROM ${BUILD_TARGET} AS final
