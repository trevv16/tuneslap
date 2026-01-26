#!/bin/bash
# E2E Test Environment Teardown Script
# Stops and removes all E2E test services and data

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "======================================"
echo "Stopping E2E Test Environment"
echo "======================================"

cd "$PROJECT_ROOT"

# Stop and remove containers, networks, and volumes
echo "Stopping Docker services and removing volumes..."
docker compose -f docker-compose.e2e.yml down -v

echo ""
echo "======================================"
echo "E2E Test Environment Cleaned Up!"
echo "======================================"
