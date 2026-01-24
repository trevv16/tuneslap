# Deployment Guide

This guide covers deploying TuneSlap using Docker Compose for self-hosting.

## Prerequisites

- Docker and Docker Compose installed
- Google Cloud Platform account (for GCS storage)
- MongoDB and Redis (included in docker-compose.yml)

## Quick Start

1. Clone the repository:

```bash
git clone https://github.com/yourusername/tuneslap.git
cd tuneslap
```

2. Create environment file:

```bash
cp server/example.env server/.env
```

3. Configure environment variables in `server/.env` (see Configuration section below)

4. Start services:

```bash
docker-compose up -d
```

**Development URLs** (localhost):

- Frontend: <http://localhost:3001>
- Backend API: <http://localhost:8082>
- API Explorer: <http://localhost:8081>
- MinIO Storage: <http://localhost:9000>

**Production URLs** (with custom domains):

- Frontend: <https://tuneslap.com>
- Backend API: <https://api.tuneslap.com>
- API Explorer: <https://openapi.tuneslap.com> (browse and test all API endpoints)
- MinIO Storage: <https://media.tuneslap.com>

## Configuration

For detailed storage configuration (S3, GCS, MinIO), see the [Storage documentation](media/storage.md).

### Environment Variables

Create a `.env` file in the `server/` directory with the following variables:

#### Required Variables

```bash
# Server Configuration
PORT=8082
JWT_SECRET=your-secret-key-here
DATABASE=tuneslap

# Database
MONGODB_URI=mongodb://mongodb:27017
REDIS_URL=redis:6379

# Google Cloud Storage
USER_UPLOADS_BUCKET=your-uploads-bucket
MEDIA_BUCKET=your-media-bucket
GOOGLE_SERVICE_ACCOUNT_EMAIL=your-service-account@project.iam.gserviceaccount.com
GOOGLE_PRIVATE_KEY_PATH=/app/keys/your-service-account-key.json
```

#### Optional Variables

```bash
# Storage Limits (in bytes, -1 for unlimited)
MAX_STORAGE_BYTES=-1

# Frontend URLs (for CORS and redirects)
# Development:
NEXT_PUBLIC_API_URL=http://localhost:8082/api/v1
NEXT_PUBLIC_SITE_URL=http://localhost:3001
# Production:
# NEXT_PUBLIC_API_URL=https://api.tuneslap.com/api/v1
# NEXT_PUBLIC_SITE_URL=https://tuneslap.com

# S3/MinIO external endpoint (for media URLs stored in database)
# Development:
S3_EXTERNAL_ENDPOINT=http://localhost:9000
# Production:
# S3_EXTERNAL_ENDPOINT=https://media.tuneslap.com
```

### Google Cloud Storage Setup

1. Create a Google Cloud Project (if you don't have one)
2. Enable the Cloud Storage API
3. Create two GCS buckets:
   - One for user uploads
   - One for processed media
4. Create a service account with Storage Admin permissions
5. Download the service account key JSON file
6. Place the key file in `server/keys/` directory
7. Update `GOOGLE_PRIVATE_KEY_PATH` to point to your key file

The `setup-keys` service in docker-compose.yml will help generate keys if needed, but you'll need to have `gcloud` CLI authenticated on your host machine.

## Docker Compose Configuration

The `docker-compose.yml` file includes:

- **server**: Go backend API
- **frontend**: Next.js frontend
- **mongodb**: MongoDB database
- **redis**: Redis cache
- **minio**: S3-compatible object storage (local development)
- **createbuckets**: Helper to initialize MinIO buckets
- **openapi**: API documentation explorer
- **setup-keys**: Optional service for GCS key setup

### Customizing Ports

To change the default ports, modify the `ports` section in `docker-compose.yml`:

```yaml
services:
  server:
    ports:
      - "8082:8082"  # Change first number to change host port
  frontend:
    ports:
      - "3001:3001"  # Change first number to change host port
```

### Persistent Storage

MongoDB and Redis data are persisted in Docker volumes:

- `mongodb_data`: MongoDB database files
- `redis_data`: Redis data

These volumes persist even when containers are stopped.

## Production Deployment

For production deployments, consider:

1. **Use a reverse proxy** (nginx, Traefik, Caddy) in front of the services
2. **Set up SSL/TLS** certificates (Let's Encrypt recommended)
3. **Configure proper CORS** settings for your domain
4. **Use environment-specific configuration** (separate .env files)
5. **Set up monitoring** and logging
6. **Configure backups** for MongoDB
7. **Use managed databases** (MongoDB Atlas, Redis Cloud) instead of containers for production

### Domain Configuration

The default `docker-compose.yml` is configured for production domains:

| Service | Domain | Internal Port |
|---------|--------|---------------|
| Frontend | tuneslap.com | 3001 |
| API Server | api.tuneslap.com | 8082 |
| MinIO Storage | media.tuneslap.com | 9000 |
| OpenAPI Docs | openapi.tuneslap.com | 8081 |

Configure your reverse proxy (Dokploy, Traefik, nginx, etc.) to route each domain to the appropriate port.

**Important**: The `S3_EXTERNAL_ENDPOINT` variable controls what URL is stored in the database for media files. Set this to your MinIO domain (e.g., `https://media.tuneslap.com`) so media URLs are accessible from the browser.

For local development, use `docker-compose.dev.yml` which defaults to localhost URLs.

### Example with Nginx

```nginx
# Frontend
server {
    listen 443 ssl;
    server_name tuneslap.com;

    location / {
        proxy_pass http://localhost:3001;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

# API Server
server {
    listen 443 ssl;
    server_name api.tuneslap.com;

    location / {
        proxy_pass http://localhost:8082;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

# MinIO Storage
server {
    listen 443 ssl;
    server_name media.tuneslap.com;

    location / {
        proxy_pass http://localhost:9000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

# OpenAPI Docs
server {
    listen 443 ssl;
    server_name openapi.tuneslap.com;

    location / {
        proxy_pass http://localhost:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Updating

To update to a new version:

```bash
# Pull latest changes
git pull

# Rebuild and restart services
docker-compose up -d --build
```

## Troubleshooting

### Services won't start

- Check that ports 3001, 8082, 27017, and 6379 are not in use
- Verify environment variables are set correctly
- Check logs: `docker-compose logs`

### Database connection issues

- Ensure MongoDB container is running: `docker-compose ps`
- Check MongoDB URI matches the service name: `mongodb://mongodb:27017`
- Verify network connectivity between containers

### Storage issues

- Ensure GCS buckets exist and are accessible
- Verify service account key file path is correct
- Check service account has proper permissions

### Viewing Logs

```bash
# All services
docker-compose logs

# Specific service
docker-compose logs server
docker-compose logs frontend

# Follow logs
docker-compose logs -f
```

## Backup and Restore

### MongoDB Backup

```bash
docker-compose exec mongodb mongodump --out /data/backup
docker cp tuneslap_mongodb_1:/data/backup ./backup
```

### MongoDB Restore

```bash
docker cp ./backup tuneslap_mongodb_1:/data/backup
docker-compose exec mongodb mongorestore /data/backup
```

## Security Considerations

- Never commit `.env` files to version control
- Use strong `JWT_SECRET` values
- Restrict GCS bucket permissions to minimum required
- Use firewall rules to limit database access
- Keep Docker images updated
- Regularly rotate service account keys
