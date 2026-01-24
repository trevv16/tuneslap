# Start from golang 1.24 base image
FROM golang:1.25-bookworm AS builder

# Install build dependencies (required for bimg/libvips) and code generation tools
# Install Node.js 20 from NodeSource for openapi-generator compatibility
RUN apt-get update && apt-get install -y --no-install-recommends \
    libvips-dev \
    gcc \
    pkg-config \
    default-jre-headless \
    curl \
    && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/*

# Work from project root structure so generate-types.sh paths work
WORKDIR /workspace

# Copy go mod and sum files first for caching
COPY server/go.mod server/go.sum ./server/

# Download all dependencies
RUN cd server && go mod download

# Copy OpenAPI specs and scripts for code generation
COPY openapi ./openapi
COPY scripts ./scripts

# Copy the server source code
COPY server ./server

# Generate Go types from OpenAPI spec using the project script
RUN npm install -g @openapitools/openapi-generator-cli swagger-cli && \
    ./scripts/generate-types.sh --backend-only

# Build the application from server directory
# CGO_ENABLED=1 is required for bimg
RUN cd server && CGO_ENABLED=1 GOOS=linux go build -o main .

# Final stage
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    libvips \
    ffmpeg \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /workspace/server/main .

# Set environment variables
# We do NOT set GO_ENV=production here because the application logic 
# only loads .env if GO_ENV is empty or "local".
ENV PORT=8082

# Expose the port
EXPOSE 8082

# Run the executable
CMD ["./main"]
