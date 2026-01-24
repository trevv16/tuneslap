.PHONY: dev fmt fmt-check swagger docs docs-stop docs-build docs-logs generate-types build test lint clean download-assets setup-keys docker-dev docker-dev-stop

dev:
	air

fmt:
	go list ./... | grep -v /generated | xargs go fmt

fmt-check:
	@if [ "$$(gofmt -s -l . | grep -v '^generated/' | wc -l)" -gt 0 ]; then \
		echo "Code is not formatted. Run 'make fmt'"; \
		gofmt -s -l . | grep -v "^generated/"; \
		exit 1; \
	fi

swagger:
	swag init --dir ./,./handlers

# Build the binary
build:
	go build -o bin/server main.go

# Run tests with race detection
test:
	go test -v -race ./...

# Run linter (uses standard library go vet)
lint:
	go vet ./...

# Clean build artifacts
clean:
	rm -rf bin

# Download demo audio assets
download-assets:
	cd .. && sh scripts/download-demo-audio.sh

# Setup service account keys (defaults to staging)
setup-keys:
	cd .. && sh scripts/setup-service-account.sh staging

# Docker dev helpers (running from root)
docker-dev:
	cd .. && docker compose -f docker-compose.dev.yml up -d

docker-dev-stop:
	cd .. && docker compose -f docker-compose.dev.yml down

docs:
	cd ../openapi && docker compose up -d

docs-stop:
	cd ../openapi && docker compose down

docs-build:
	cd ../openapi && docker compose build

docs-logs:
	cd ../openapi && docker compose logs -f

generate-types:
	cd .. && sh scripts/generate-types.sh
