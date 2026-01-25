#!/bin/bash
# E2E Test Environment Setup Script
# Starts all required services for running E2E tests locally

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "======================================"
echo "Starting E2E Test Environment"
echo "======================================"

cd "$PROJECT_ROOT"

# Check if docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Error: Docker is not running. Please start Docker and try again."
    exit 1
fi

# Start services with docker compose
echo "Starting Docker services..."
docker compose -f docker-compose.e2e.yml up -d --wait

# Wait for backend health check
echo "Waiting for backend to be healthy..."
MAX_RETRIES=30
RETRY_DELAY=2
for i in $(seq 1 $MAX_RETRIES); do
    if curl -sf http://localhost:8082/health > /dev/null 2>&1; then
        echo "Backend is healthy!"
        break
    fi
    if [ $i -eq $MAX_RETRIES ]; then
        echo "Error: Backend did not become healthy in time"
        echo "Checking container logs..."
        docker compose -f docker-compose.e2e.yml logs server
        exit 1
    fi
    echo "Waiting for backend... ($i/$MAX_RETRIES)"
    sleep $RETRY_DELAY
done

echo ""
echo "======================================"
echo "E2E Test Environment Ready!"
echo "======================================"
echo ""
echo "Services running:"
echo "  - Backend API:  http://localhost:8082"
echo "  - MongoDB:      mongodb://localhost:27017"
echo "  - Redis:        localhost:6379"
echo "  - MinIO:        http://localhost:9000"
echo "  - MinIO Console: http://localhost:9001"
echo ""
echo "To run E2E tests:"
echo "  cd frontend && yarn test:e2e"
echo ""
echo "To stop the environment:"
echo "  ./scripts/e2e-teardown.sh"
echo ""
