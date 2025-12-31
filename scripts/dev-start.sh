#!/bin/bash

# Development startup script for Streamly
# Starts the development environment with Docker Compose

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
COMPOSE_FILE="deployments/docker-compose.dev.yml"
PROJECT_NAME="streamly-dev"

print_usage() {
    echo -e "${BLUE}Usage: $0 [command]${NC}"
    echo
    echo -e "${YELLOW}Commands:${NC}"
    echo "  start     Start all development services (default)"
    echo "  stop      Stop all services"
    echo "  restart   Restart all services"
    echo "  logs      Show logs for all services"
    echo "  status    Show status of all services"
    echo "  clean     Stop and remove all containers, networks, and volumes"
    echo "  frontend  Start only frontend service"
    echo "  backend   Start only backend service"
    echo
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0              # Start all services"
    echo "  $0 start        # Start all services"
    echo "  $0 logs         # Show logs"
    echo "  $0 frontend     # Start only frontend"
}

check_requirements() {
    echo -e "${BLUE}🔍 Checking requirements...${NC}"

    if ! command -v docker &> /dev/null; then
        echo -e "${RED}❌ Docker is not installed${NC}"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        echo -e "${RED}❌ Docker Compose is not installed${NC}"
        exit 1
    fi

    if [[ ! -f "$COMPOSE_FILE" ]]; then
        echo -e "${RED}❌ Docker Compose file not found: $COMPOSE_FILE${NC}"
        exit 1
    fi

    echo -e "${GREEN}✅ Requirements check passed${NC}"
}

start_services() {
    local service=${1:-""}

    echo -e "${YELLOW}🚀 Starting development environment...${NC}"

    if [[ -n "$service" ]]; then
        echo -e "${BLUE}Service: $service${NC}"
        docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" up -d "$service"
    else
        echo -e "${BLUE}Starting all services${NC}"
        docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" up -d
    fi

    echo -e "${GREEN}✅ Services started successfully${NC}"
    show_status
}

stop_services() {
    echo -e "${YELLOW}⏹️  Stopping development environment...${NC}"
    docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" down
    echo -e "${GREEN}✅ Services stopped${NC}"
}

restart_services() {
    echo -e "${YELLOW}🔄 Restarting development environment...${NC}"
    stop_services
    sleep 2
    start_services
}

show_logs() {
    echo -e "${BLUE}📋 Showing logs...${NC}"
    docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" logs -f
}

show_status() {
    echo -e "${BLUE}📊 Service Status:${NC}"
    docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" ps
    echo
    echo -e "${BLUE}🌐 Access URLs:${NC}"
    echo -e "Frontend:       ${GREEN}http://localhost:3000${NC}"
    echo -e "Platform:       ${GREEN}http://localhost:8094${NC}"
    echo -e "Ingestor:       ${GREEN}http://localhost:8091${NC}"
    echo -e "Query API:      ${GREEN}http://localhost:8092${NC}"
    echo -e "Auth:           ${GREEN}http://localhost:8093${NC}"
    echo
    echo -e "${BLUE}🔐 Zitadel (Auth):${NC}"
    echo -e "Console:        ${GREEN}http://localhost:8080/ui/console${NC}"
    echo -e "Credentials:    ${YELLOW}root@streamly.localhost / RootPassword123!${NC}"
    echo -e "${YELLOW}Note: Using built-in Zitadel login UI (localhost issue fixed)${NC}"
    echo
    echo -e "${BLUE}🛠️  Admin Tools:${NC}"
    echo -e "Adminer (DB):   ${GREEN}http://localhost:8081${NC}"
    echo -e "ClickHouse UI:  ${GREEN}http://localhost:8124${NC}"
    echo -e "MailHog:        ${GREEN}http://localhost:8025${NC}"
}

clean_environment() {
    echo -e "${YELLOW}🧹 Cleaning development environment...${NC}"
    echo -e "${RED}⚠️  This will remove all containers, networks, and volumes!${NC}"
    read -p "Are you sure? (y/N): " -n 1 -r
    echo

    if [[ $REPLY =~ ^[Yy]$ ]]; then
        docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" down -v --remove-orphans
        docker system prune -f
        echo -e "${GREEN}✅ Environment cleaned${NC}"
    else
        echo -e "${BLUE}ℹ️  Clean operation cancelled${NC}"
    fi
}

main() {
    local command=${1:-start}

    if [[ "$command" == "-h" || "$command" == "--help" ]]; then
        print_usage
        exit 0
    fi

    cd "$(dirname "$0")/.."

    check_requirements

    case $command in
        "start")
            start_services
            ;;
        "stop")
            stop_services
            ;;
        "restart")
            restart_services
            ;;
        "logs")
            show_logs
            ;;
        "status")
            show_status
            ;;
        "clean")
            clean_environment
            ;;
        "frontend")
            start_services "frontend"
            ;;
        "backend")
            start_services "backend"
            ;;
        *)
            echo -e "${RED}❌ Unknown command: $command${NC}"
            print_usage
            exit 1
            ;;
    esac
}

# Run main function
main "$@"

