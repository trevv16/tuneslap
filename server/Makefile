dev:
	air

fmt:
	go list ./... | grep -v /generated | xargs go fmt

swagger:
	swag init --dir ./,./handlers

docs:
	cd ../openapi && docker compose up -d

docs-stop:
	cd ../openapi && docker compose down

docs-build:
	cd ../openapi && docker compose build

docs-logs:
	cd ../openapi && docker compose logs -f

.PHONY: generate-types
generate-types:
	cd .. && sh scripts/generate-types.sh