# Start from golang 1.24 base image
FROM golang:1.25-bookworm AS builder

# Install build dependencies (required for bimg/libvips) and code generation tools
# Install Node.js 24 from NodeSource for openapi-generator
# Package versions pinned for reproducibility (Debian bookworm)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libvips-dev=8.14.1-3 \
    gcc=4:12.2.0-3 \
    pkg-config=1.8.1-1 \
    default-jre-headless=2:1.17-74 \
    curl=7.88.1-10+deb12u8 \
    && curl -fsSL https://deb.nodesource.com/setup_24.x | bash - \
    && apt-get install -y --no-install-recommends nodejs \
    && rm -rf /var/lib/apt/lists/*

# Work from project root structure so generate-types.sh paths work
WORKDIR /workspace

# Copy go mod and sum files first for caching
COPY server/go.mod server/go.sum ./server/

# Download all dependencies
WORKDIR /workspace/server
RUN go mod download

# Copy OpenAPI specs and scripts for code generation
WORKDIR /workspace
COPY openapi ./openapi
COPY scripts ./scripts

# Copy the server source code
COPY server ./server

# Generate Go types from OpenAPI spec using the project script
RUN npm install -g @openapitools/openapi-generator-cli swagger-cli && \
    ./scripts/generate-types.sh --backend-only

# Build the application from server directory
# CGO_ENABLED=1 is required for bimg
WORKDIR /workspace/server
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Final stage
FROM debian:bookworm-slim

# Install runtime dependencies
# Package versions pinned for reproducibility (Debian bookworm-slim)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libvips42=8.14.1-3 \
    ffmpeg=7:5.1.6-0+deb12u1 \
    ca-certificates=20230311 \
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
