# Upload Flow

This document describes how files are uploaded to TuneSlap using signed URLs and direct-to-storage transfers.

## How It Works

TuneSlap uses a three-step upload process:

1. **Generate Signed URL** - Backend creates a time-limited URL for direct upload.
2. **Direct Upload** - Frontend uploads the file directly to object storage.
3. **Create Media Record** - Backend creates the media record and queues processing.

This approach keeps large files off the backend server and allows for progress tracking on the frontend.

## Step 1: Generate Signed URL

The frontend requests a signed URL by calling `POST /api/v1/media/upload-url`.

**Request:**

```json
{
  "fileName": "my-sound-effect.mp3",
  "contentType": "audio/mpeg",
  "fileSize": 1048576
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "signedUrl": "https://storage.googleapis.com/...",
    "objectKey": "user123/audio/my-sound-effect.mp3",
    "bucketName": "tuneslap-user-uploads"
  }
}
```

The backend validates:
- File type is allowed (audio or image)
- File size is within storage limits
- User is authenticated

## Step 2: Direct Upload

The frontend uploads the file directly to object storage using the signed URL:

```javascript
const response = await fetch(signedUrl, {
  method: 'PUT',
  headers: {
    'Content-Type': contentType,
  },
  body: file,
});
```

The signed URL is valid for 15 minutes and only allows PUT requests with the specified content type.

## Step 3: Create Media Record

After a successful upload, the frontend creates the media record:

**Request to `POST /api/v1/media`:**

```json
{
  "fileName": "my-sound-effect.mp3",
  "mediaType": "audio",
  "fileUrl": "https://storage.googleapis.com/bucket/user123/audio/my-sound-effect.mp3",
  "fileSize": 1048576,
  "contentType": "audio/mpeg",
  "processingParams": {
    "audio": {
      "contentType": "audio/webm",
      "normalize": true
    }
  }
}
```

The backend:
1. Validates the request
2. Creates a media record with status `pending`
3. Queues a background task to process the file
4. Returns the created media record

## Object Key Structure

Files are organized in storage using this pattern:

```
{userId}/{mediaType}/{fileName}
```

For example:
- `abc123/audio/intro-music.mp3`
- `abc123/image/logo.png`

## Security

- **Signed URLs** expire after 15 minutes
- **Content-Type validation** ensures the uploaded file matches the expected type
- **Storage limits** are checked before generating the URL
- **JWT authentication** is required for all endpoints

## Error Handling

| Error | Cause | Solution |
|-------|-------|----------|
| 401 Unauthorized | Missing or invalid JWT | Re-authenticate |
| 400 Bad Request | Invalid file type or size | Check file before upload |
| 413 Payload Too Large | File exceeds storage limit | Reduce file size or upgrade plan |
| 500 Internal Server Error | Storage service unavailable | Retry later |

## Related

- [Storage](./storage.md) - Storage provider configuration
- [Processing](./processing.md) - What happens after upload
- [Async Operations](./async-operations.md) - Background task queue
