# streamly
Streamly.io SaaS Platform to monitor your event service fully transparently

## Development Setup

### Prerequisites
- Docker & Docker Compose
- Go 1.25+
- Bun (for frontend development)

### Quick Start

```bash
# Start all services
cd deployments
docker-compose -f docker-compose.dev.yml up

# Services will be available at:
# - Frontend: http://localhost:3000
# - Ingestor Service: http://localhost:8091
# - Query Service: http://localhost:8092
# - Platform Service: http://localhost:8084
# - Auth Service: http://localhost:8093
# - Zitadel: http://localhost:8888
# - Zitadel Login: http://localhost:3003
```

### Hot Reload Development

All Go services (platform, ingestion, query) are configured with [Air](https://github.com/cosmtrek/air) for automatic hot reload during development.

**How it works:**
- Each service has an `.air.toml` configuration in its `cmd/` directory
- Code changes are detected automatically and trigger rebuilds
- Services restart in **< 5 seconds** after code changes
- Volume mounts ensure your local code is synced with containers

**Verify hot reload:**
1. Start the services: `docker-compose -f deployments/docker-compose.dev.yml up`
2. Make a code change in any Go service (e.g., `cmd/platform/main.go`)
3. Watch the service logs - you'll see Air rebuild and restart the service
4. Changes are live in < 5 seconds

**Development workflow:**
```bash
# Watch logs for a specific service
docker-compose -f deployments/docker-compose.dev.yml logs -f platform-service

# Restart a service manually (if needed)
docker-compose -f deployments/docker-compose.dev.yml restart platform-service

# Rebuild a service (if dependencies changed)
docker-compose -f deployments/docker-compose.dev.yml up -d --build platform-service
```

### Project Structure

```
streamly/
├── cmd/                    # Service entry points
│   ├── platform/          # Platform service (orgs, users, tenants)
│   ├── ingestion/         # Event ingestion service
│   ├── query-api/         # Query API service
│   └── auth/              # Authentication service
├── internal/              # Internal packages
│   ├── core/             # Business logic
│   └── storage/          # Database adapters
├── deployments/           # Docker compose & configs
│   ├── docker-compose.dev.yml
│   └── zitadel/          # Auth configs
├── scripts/               # Database migrations & scripts
│   └── sql/              # PostgreSQL migrations
└── frontend/              # React frontend (Bun + TanStack)
```

### Tech Stack

**Backend:**
- Go 1.24 with hot reload (Air)
- PostgreSQL (metadata, users, orgs)
- ClickHouse (event analytics)
- Redis (caching, pub/sub)
- RabbitMQ (event queue)

**Auth:**
- Zitadel (OAuth/OIDC provider)
- Service accounts for inter-service auth

**Frontend:**
- React + TypeScript
- Bun (runtime & package manager)
- TanStack Router + Query
- Shadcn UI components
