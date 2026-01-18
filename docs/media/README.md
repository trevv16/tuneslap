# Media System

This section documents the media handling system in TuneSlap, covering file uploads, storage, processing, and background tasks.

## Contents

- [Upload Flow](./upload-flow.md) - How files are uploaded using signed URLs and direct-to-storage transfers.
- [Storage](./storage.md) - Object storage interface supporting GCS and S3/MinIO.
- [Processing](./processing.md) - Audio and image processing parameters and pipelines.
- [Async Operations](./async-operations.md) - Background task queue for media processing.

## Overview

TuneSlap uses a two-bucket storage architecture:

1. **User Uploads Bucket** - Receives raw files uploaded directly from the frontend via signed URLs.
2. **Media Bucket** - Stores processed/normalized files ready for playback.

When a user uploads a file:

1. Frontend requests a signed upload URL from the backend.
2. Frontend uploads the file directly to the User Uploads bucket.
3. Frontend confirms the upload by creating a media record.
4. Backend queues a background task to process the file.
5. The processed file is stored in the Media bucket and the record is updated.

```
┌──────────┐     1. Request URL      ┌──────────┐
│ Frontend │ ────────────────────────▶│ Backend  │
└──────────┘                          └──────────┘
     │                                      │
     │  2. Upload file                      │
     ▼                                      │
┌──────────────────┐                        │
│ User Uploads     │                        │
│ Bucket           │                        │
└──────────────────┘                        │
     │                                      │
     │  3. Confirm upload                   │
     └─────────────────────────────────────▶│
                                            │
                                            │ 4. Queue task
                                            ▼
                                      ┌──────────┐
                                      │  Redis   │
                                      │  Queue   │
                                      └──────────┘
                                            │
                                            │ 5. Process
                                            ▼
                                      ┌──────────────────┐
                                      │ Media Bucket     │
                                      └──────────────────┘
```
