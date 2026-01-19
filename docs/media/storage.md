# Storage

TuneSlap uses a pluggable storage interface that supports multiple providers. The default is S3-compatible storage (including MinIO for local development), with Google Cloud Storage available for production.

## Providers

### S3-Compatible (Default)

Works with AWS S3, MinIO, Cloudflare R2, DigitalOcean Spaces, and other S3-compatible services.

**Environment Variables:**

| Variable | Required | Description |
|----------|----------|-------------|
| `STORAGE_PROVIDER` | No | Set to `s3` (default if not set) |
| `S3_ENDPOINT` | Yes* | Endpoint URL (e.g., `http://minio:9000`). Required for non-AWS S3. |
| `S3_EXTERNAL_ENDPOINT` | No | Public endpoint for media URLs. See [External Endpoint Configuration](#external-endpoint-configuration) below. |
| `S3_REGION` | No | AWS region (defaults to `us-east-1`) |
| `S3_ACCESS_KEY` | Yes | Access key ID |
| `S3_SECRET_KEY` | Yes | Secret access key |
| `USER_UPLOADS_BUCKET` | Yes | Bucket for raw user uploads |
| `MEDIA_BUCKET` | Yes | Bucket for processed media |

*Not required for AWS S3.

### Google Cloud Storage

For production deployments on Google Cloud.

**Environment Variables:**

| Variable | Required | Description |
|----------|----------|-------------|
| `STORAGE_PROVIDER` | Yes | Set to `gcs` |
| `GOOGLE_APPLICATION_CREDENTIALS` | No | Path to service account JSON (uses ADC if not set) |
| `GOOGLE_PRIVATE_KEY_PATH` | Yes | Path to service account key file for signed URLs |
| `GOOGLE_SERVICE_ACCOUNT_EMAIL` | No | Service account email (read from key file if not set) |
| `USER_UPLOADS_BUCKET` | Yes | Bucket for raw user uploads |
| `MEDIA_BUCKET` | Yes | Bucket for processed media |

## Two-Bucket Architecture

TuneSlap uses two separate buckets:

1. **User Uploads Bucket** (`USER_UPLOADS_BUCKET`)
   - Receives raw files uploaded directly from the frontend
   - Files are temporary and can be deleted after processing
   - Configured with short retention if desired

2. **Media Bucket** (`MEDIA_BUCKET`)
   - Stores processed/normalized files
   - Served to users for playback
   - Should have longer retention and CDN if needed

## ObjectStorage Interface

All storage providers implement this interface (defined in `server/services/storage/storage.go`):

```go
type ObjectStorage interface {
    UploadFile(ctx context.Context, req UploadFileRequest) error
    DownloadFile(ctx context.Context, objectName, destPath string) error
    DeleteFile(ctx context.Context, objectName string) error
    GenerateSignedUploadURL(ctx context.Context, objectName, contentType string, expires time.Duration) (string, error)
    GenerateSignedDownloadURL(ctx context.Context, objectName string, expires time.Duration) (string, error)
    GetBucketName() string
    GetFileURL(objectName string) string
}
```

## Object Key Structure

Files are stored using this key pattern:

```
{userId}/{mediaType}/{fileName}
```

Examples:
- `abc123/audio/intro-music.mp3`
- `abc123/image/logo.png`

The helper function `storage.GetMediaKey()` generates these keys.

## External Endpoint Configuration

The `S3_EXTERNAL_ENDPOINT` variable is critical for production deployments. It controls the URL that gets stored in the database for media files and must be accessible from the browser.

**How it works:**

1. `S3_ENDPOINT` is used internally by the server to communicate with MinIO (e.g., `http://minio:9000`)
2. `S3_EXTERNAL_ENDPOINT` is used to build public URLs stored in the database (e.g., `https://media.tuneslap.com`)

**Example URLs generated:**

| Environment | S3_EXTERNAL_ENDPOINT | Stored URL |
|-------------|---------------------|------------|
| Development | `http://localhost:9000` | `http://localhost:9000/tuneslap-media/userId/audio/file.mp3` |
| Production | `https://media.tuneslap.com` | `https://media.tuneslap.com/tuneslap-media/userId/audio/file.mp3` |

**Configuration examples:**

```bash
# Development (localhost)
S3_ENDPOINT=http://minio:9000
S3_EXTERNAL_ENDPOINT=http://localhost:9000

# Production (custom domain)
S3_ENDPOINT=http://minio:9000
S3_EXTERNAL_ENDPOINT=https://media.tuneslap.com
```

For production, configure your reverse proxy to route `media.tuneslap.com` to MinIO port 9000.

## Local Development with MinIO

Use `docker-compose.dev.yml` for local development, which defaults to localhost URLs:

```yaml
minio:
  image: minio/minio:latest
  command: server /data --console-address ":9001"
  ports:
    - "9000:9000"   # API
    - "9001:9001"   # Web console
  environment:
    - MINIO_ROOT_USER=minioadmin
    - MINIO_ROOT_PASSWORD=minioadmin
```

Access the MinIO console at `http://localhost:9001` with credentials `minioadmin/minioadmin`.

Buckets are created automatically by the `createbuckets` service on startup.

## CORS Configuration

For direct uploads from the browser, buckets need CORS configured to allow PUT requests.

### MinIO

MinIO allows all origins by default in development. For production, configure via `mc`:

```bash
mc anonymous set public myminio/tuneslap-uploads
```

### GCS

Update bucket CORS via `gsutil`:

```bash
gsutil cors set cors.json gs://your-bucket-name
```

Example `cors.json`:

```json
[
  {
    "origin": ["https://tuneslap.com"],
    "method": ["GET", "PUT", "POST"],
    "responseHeader": ["Content-Type"],
    "maxAgeSeconds": 3600
  }
]
```

## Related

- [Upload Flow](./upload-flow.md) - How uploads use storage
- [Processing](./processing.md) - How processed files are stored
