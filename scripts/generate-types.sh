#!/bin/bash
set -e

GENERATE_FRONTEND=true
GENERATE_BACKEND=true

# Parse command line arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --frontend-only) GENERATE_BACKEND=false ;;
        --backend-only) GENERATE_FRONTEND=false ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

if [ "$GENERATE_FRONTEND" = false ] && [ "$GENERATE_BACKEND" = false ]; then
    echo "Error: Cannot specify both --frontend-only and --backend-only"
    exit 1
fi

echo "🚀 Bundling OpenAPI spec..."
npx swagger-cli bundle --dereference openapi/openapi.yaml -o openapi/bundled.json

if [ "$GENERATE_BACKEND" = true ]; then
    echo "🚀 Generating Go types from OpenAPI spec..."

    # Ensure the output directory exists
    mkdir -p server/generated

    # Clean up old generated files before regenerating
    # Remove all generated model files (they'll be regenerated)
    find server/generated -maxdepth 1 -type f -name "model_*.go" -delete 2>/dev/null || true
    find server/generated -maxdepth 1 -type f -name "api_*.go" -delete 2>/dev/null || true
    # Keep docs and other directories for now, but remove old response docs
    find server/generated/docs -name "*200*.md" -o -name "*400*.md" -o -name "*201*.md" 2>/dev/null | xargs rm -f 2>/dev/null || true

    # Generate Go types using openapi-generator-cli
    # Run from project root, output to server/generated
    npx @openapitools/openapi-generator-cli generate \
      -i openapi/bundled.json \
      -g go \
      -o server/generated \
      --type-mappings=string:objectid=primitive.ObjectID \
      --import-mappings=ObjectID=go.mongodb.org/mongo-driver/bson/primitive \
      --additional-properties=packageName=api,packageVersion=1.0.0,generateInterfaces=true,enumClassPrefix=true,generateTests=false

    # Remove generated test files and go.mod (we use the main module's go.mod)
    rm -rf server/generated/test
    rm -f server/generated/go.mod server/generated/go.sum

    # Clean up duplicate response files with status codes (keep the ones with clean names)
    find server/generated -maxdepth 1 -type f -name "model_*_200_response.go" -delete 2>/dev/null || true
    find server/generated -maxdepth 1 -type f -name "model_*_201_response.go" -delete 2>/dev/null || true
    find server/generated -maxdepth 1 -type f -name "model_*_400_response.go" -delete 2>/dev/null || true
fi

if [ "$GENERATE_FRONTEND" = true ]; then
    echo "🚀 Generating TypeScript API client from OpenAPI spec..."

    # Ensure the output directory exists
    mkdir -p frontend/api

    # Clean up old generated files (but preserve custom files)
    # Backup custom files
    [ -f frontend/api/config.ts ] && cp frontend/api/config.ts /tmp/api-config.ts.bak || true
    [ -f frontend/api/uploadUrl.ts ] && cp frontend/api/uploadUrl.ts /tmp/api-uploadUrl.ts.bak || true

    # Remove generated directories and files
    rm -rf frontend/api/models
    rm -rf frontend/api/apis
    rm -rf frontend/api/runtime
    # Remove generated files in root but keep custom ones
    find frontend/api -maxdepth 1 -type f \( -name "*.ts" -o -name "*.tsx" \) ! -name "config.ts" ! -name "index.ts" ! -name "uploadUrl.ts" -delete 2>/dev/null || true

    # Generate TypeScript-fetch API client
    npx @openapitools/openapi-generator-cli generate \
      -i openapi/bundled.json \
      -g typescript-fetch \
      -o frontend/api \
      --additional-properties=supportsES6=true,withInterfaces=true,useSingleRequestParameter=true,modelPropertyNaming=camelCase,fileNaming=PascalCase,stringEnums=true

    # Restore custom files if they were backed up
    [ -f /tmp/api-config.ts.bak ] && mv /tmp/api-config.ts.bak frontend/api/config.ts || true
    [ -f /tmp/api-uploadUrl.ts.bak ] && mv /tmp/api-uploadUrl.ts.bak frontend/api/uploadUrl.ts || true
fi

echo "✅ Type generation completed successfully!"
