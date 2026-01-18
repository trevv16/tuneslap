# Start from golang 1.24 base image
FROM golang:1.24-bookworm AS builder

# Install build dependencies (required for bimg/libvips)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libvips-dev \
    gcc \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=1 is required for bimg
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

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
COPY --from=builder /app/main .

# Set environment variables
# We do NOT set GO_ENV=production here because the application logic 
# only loads .env if GO_ENV is empty or "local".
ENV PORT=8080

# Expose the port
EXPOSE 8080

# Run the executable
CMD ["./main"]
