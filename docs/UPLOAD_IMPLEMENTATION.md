# File Upload Implementation

This document describes the implementation of direct file upload to Google Cloud Storage (GCS) for the TuneSlap application.

## Overview

The file upload system uses a two-step process:
1. **Generate Signed URL**: Backend generates a time-limited signed URL for direct upload to GCS
2. **Direct Upload**: Frontend uploads the file directly to GCS using the signed URL
3. **Create Media Record**: Backend creates a media record with the uploaded file URL

## Architecture

### Backend Components

#### 1. GCS Service (`server/services/gcs.go`)
- `GenerateSignedUploadURL()`: Creates signed URLs for PUT requests
- `GenerateSignedDownloadURL()`: Creates signed URLs for GET requests
- Uses GCS v4 signing scheme for security

#### 2. Media Handler (`server/handlers/media_handler.go`)
- `HandleGenerateUploadURL()`: New endpoint for generating upload URLs
- Validates file size against storage limit
- Generates object keys based on user ID, media type, and filename

#### 3. Router (`server/router/main.go`)
- Added `/api/v1/media/upload-url` endpoint
- Protected with JWT authentication

### Frontend Components

#### 1. Media API (`frontend/api/media.ts`)
- `generateUploadUrl()`: Calls backend to get signed URL
- Handles response parsing and error handling

#### 2. Media Hooks (`frontend/hooks/media.ts`)
- `useFileUpload()`: Hook for file upload operations
- `useMediaCreate()`: Updated to use file upload before creating media record

#### 3. Create Media Form (`frontend/app/(authenticated)/library/components/CreateMediaForm.tsx`)
- Updated to show upload progress
- Integrated with new upload flow

## API Endpoints

### POST `/api/v1/media/upload-url`
Generates a signed URL for direct file upload.

**Request Body:**
```json
{
  "fileName": "my-audio-file.mp3",
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
    "objectKey": "user123/audio/my-audio-file.mp3",
    "bucketName": "tuneslap-user-uploads-staging"
  }
}
```

## File Upload Flow

1. **User selects file** in the Create Media Form
2. **Frontend validates** file type, size, and storage limits
3. **Generate signed URL** by calling `/api/v1/media/upload-url`
4. **Upload file directly** to GCS using the signed URL
5. **Create media record** by calling `/api/v1/media` with the file URL
6. **Show success/error** message to user

## Security Features

- **Signed URLs**: Time-limited (15 minutes) and scoped to specific operations
- **Storage Validation**: Checks storage limits before allowing upload
- **File Type Validation**: Validates content types on both frontend and backend
- **JWT Authentication**: All endpoints require valid authentication
- **CORS Configuration**: Updated to allow PUT requests for direct uploads

## Infrastructure Updates

### CORS Configuration
Updated GCS bucket CORS settings to allow:
- `PUT` and `POST` methods for file uploads
- Additional response headers for CORS support

### Environment Variables
Required environment variables:
- `USER_UPLOADS_BUCKET`: GCS bucket for user uploads
- `MEDIA_BUCKET`: GCS bucket for processed media
- Google Cloud credentials (via service account or default credentials)

## Error Handling

### Frontend Errors
- Network errors during upload
- Invalid file types or sizes
- Storage limit exceeded
- Authentication failures

### Backend Errors
- GCS service unavailable
- Invalid request parameters
- Storage limit exceeded
- Authentication/authorization failures

## Testing

### Manual Testing
1. Upload a small image file (< 1MB)
2. Upload a larger audio file (> 10MB)
3. Test with invalid file types
4. Test with files exceeding storage limits
5. Test upload cancellation

### Automated Testing
- Unit tests for GCS service functions
- Integration tests for upload flow
- E2E tests for complete user journey

## Performance Considerations

- **Direct Upload**: Files go directly to GCS, reducing server load
- **Progress Tracking**: Real-time upload progress feedback
- **File Size Limits**: Configurable limits based on env variable
- **Concurrent Uploads**: Multiple files can be uploaded simultaneously

## Future Enhancements

1. **Resumable Uploads**: Support for large file uploads with resume capability
2. **Chunked Uploads**: Split large files into chunks for better reliability
3. **Upload Queue**: Queue system for managing multiple uploads
4. **Background Processing**: Process uploaded files in background tasks
5. **CDN Integration**: Use CDN for faster file delivery 