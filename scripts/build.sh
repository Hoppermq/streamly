#!/bin/bash

# Build script for Streamly - supports multiple environments and architectures
# Usage: ./scripts/build.sh [environment] [service] [architecture]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
ENVIRONMENT=${1:-development}
SERVICE=${2:-all}
ARCHITECTURE=${3:-amd64}

# Supported values
SUPPORTED_ENVS="development production"
SUPPORTED_SERVICES="frontend backend all"
SUPPORTED_ARCHS="amd64 arm64"

print_usage() {
    echo -e "${BLUE}Usage: $0 [environment] [service] [architecture]${NC}"
    echo
    echo -e "${YELLOW}Arguments:${NC}"
    echo "  environment: ${SUPPORTED_ENVS} (default: development)"
    echo "  service:     ${SUPPORTED_SERVICES} (default: all)"
    echo "  architecture: ${SUPPORTED_ARCHS} (default: amd64)"
    echo
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0                          # Build all services for development (amd64)"
    echo "  $0 production               # Build all services for production (amd64)"
    echo "  $0 development frontend     # Build only frontend for development"
    echo "  $0 production backend arm64 # Build backend for production on ARM64"
}

validate_args() {
    if [[ ! " $SUPPORTED_ENVS " =~ " $ENVIRONMENT " ]]; then
        echo -e "${RED}Error: Unsupported environment '$ENVIRONMENT'${NC}"
        echo -e "Supported: $SUPPORTED_ENVS"
        exit 1
    fi

    if [[ ! " $SUPPORTED_SERVICES " =~ " $SERVICE " ]]; then
        echo -e "${RED}Error: Unsupported service '$SERVICE'${NC}"
        echo -e "Supported: $SUPPORTED_SERVICES"
        exit 1
    fi

    if [[ ! " $SUPPORTED_ARCHS " =~ " $ARCHITECTURE " ]]; then
        echo -e "${RED}Error: Unsupported architecture '$ARCHITECTURE'${NC}"
        echo -e "Supported: $SUPPORTED_ARCHS"
        exit 1
    fi
}

print_info() {
    echo -e "${BLUE}🚀 Building Streamly${NC}"
    echo -e "Environment:  ${GREEN}$ENVIRONMENT${NC}"
    echo -e "Service:      ${GREEN}$SERVICE${NC}"
    echo -e "Architecture: ${GREEN}$ARCHITECTURE${NC}"
    echo
}

build_frontend() {
    echo -e "${YELLOW}📦 Building Frontend...${NC}"

    local target_stage="development"
    if [[ "$ENVIRONMENT" == "production" ]]; then
        target_stage="production"
    fi

    docker build \
        --platform linux/$ARCHITECTURE \
        --target $target_stage \
        --build-arg BUILD_TARGET=$target_stage \
        --build-arg BUN_VERSION=1.0 \
        --tag streamly-frontend:$ENVIRONMENT-$ARCHITECTURE \
        --file docker/frontend.Dockerfile \
        .

    echo -e "${GREEN}✅ Frontend build completed${NC}"
}

build_backend() {
    echo -e "${YELLOW}🔧 Building Backend...${NC}"

    local target_os="linux"
    local service_path="./cmd/server"

    docker build \
        --platform linux/$ARCHITECTURE \
        --build-arg TARGET_OS=$target_os \
        --build-arg TARGET_ARCH=$ARCHITECTURE \
        --build-arg SERVICE_PATH=$service_path \
        --tag streamly-backend:$ENVIRONMENT-$ARCHITECTURE \
        --file docker/services.Dockerfile \
        .

    echo -e "${GREEN}✅ Backend build completed${NC}"
}

build_all() {
    echo -e "${YELLOW}🏗️  Building All Services...${NC}"
    build_frontend
    build_backend
    echo -e "${GREEN}✅ All services built successfully${NC}"
}

cleanup_build_cache() {
    echo -e "${YELLOW}🧹 Cleaning up build cache...${NC}"
    docker builder prune -f
    echo -e "${GREEN}✅ Build cache cleaned${NC}"
}

main() {
    if [[ "$1" == "-h" || "$1" == "--help" ]]; then
        print_usage
        exit 0
    fi

    validate_args
    print_info

    cd "$(dirname "$0")/.."

    case $SERVICE in
        "frontend")
            build_frontend
            ;;
        "backend")
            build_backend
            ;;
        "all")
            build_all
            ;;
    esac

    # Optional cleanup (uncomment if desired)
    # cleanup_build_cache

    echo
    echo -e "${GREEN}🎉 Build completed successfully!${NC}"
    echo -e "${BLUE}Tagged images:${NC}"
    docker images | grep streamly | grep $ENVIRONMENT-$ARCHITECTURE
}

main "$@"
