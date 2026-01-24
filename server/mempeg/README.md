# 🌍 **Mempeg – Media Manipulation Service**

A backend and frontend system for handling real-time and batch media processing (audio + image) with powerful features and smooth user experience.

---

### 🔧 **Core Use Case**

Users upload audio or image media for their soundboard app. The app allows:

- **Previews** of media edits on the front end
- **Permanent transformations** on the back end
- Tracking of **original and processed file URLs**
- Monitoring of **processing activity and job statuses**
- Format support for modern delivery (WebM, WebP, MP3, PNG, etc.)

---

### 🎛 **Features**

### 🟣 **Audio Processing**

- Transcode to WebM and MP3
- Trim (in/out)
- Normalize loudness
- Extract metadata (duration, sample rate, etc.)
- Advanced edits:
  - Speed adjustment
  - Pitch shifting
  - Fade in/out
  - Loop prep (auto trim silence, align start/end)
- Generate waveform preview images

### 🟢 **Image Processing**

- Resize
- Compress
- Convert to WebP and PNG
- Generate image previews (thumbnails)
- Crop
- Extract EXIF metadata
- Strip EXIF for privacy
- Auto aspect correction
- Remove background
- Apply filters (grayscale, blur, etc.)

---

### ⚙️ **Media Lifecycle & Flow (Go Backend)**

1. **Upload**
    - Frontend handles previews and UX.
    - User uploads to `/api/media/upload`.
    - File is saved to S3 (original form).
    - Basic metadata is stored in MongoDB (`Media` collection).
2. **Processing**
    - Go server pushes a job to Redis (via Asynq).
    - Worker fetches file from S3, processes using FFmpeg or `bimg` (libvips).
    - Processed file is uploaded to S3 (new key).
    - MongoDB entry is updated: `status`, `processedUrl`, `processingMetadata`.
3. **API Endpoints**
    - `POST /api/media/upload` – Upload media file.
    - `POST /api/media/:id/process` – Trigger processing with params.
    - `GET /api/media/:id` – Get metadata & status.
    - Future:
        - `PATCH /api/media/:id/reprocess` – Reprocess with new parameters.
        - `DELETE /api/media/:id` – Delete file & metadata.

---

### 🧠 **Tracked Metadata (MongoDB)**

- `fileName` (string)
- `fileSize` (number, bytes)
- `contentType` (string, e.g., `audio/webm`)
- `mediaType` (`audio` | `image`)
- `originalUrl` (string)
- `processedUrl` (string)
- `processingStatus` (`pending`, `processing`, `done`, `error`)
- `processingParams` (object: trim, fade, speed, etc.)
- `duration` (number, seconds)
- `dimensions` (object: width, height)
- `codec`, `bitrate`, `waveformUrl`, `thumbnailUrl` (optional, by type)
- `createdAt`, `updatedAt`, `userId`

---

### 🏗 **Architecture**

### 🟨 Frontend

- **Remix** + **Vite**
- **TanStack Query** for caching/data fetching
- **Tailwind CSS** for UI
- **Web Audio API** + **Waveform API** for client-side previews
- Real-time preview edits before submitting to the backend

### 🟦 Backend (Go Only)

- **Go + Fiber** – REST API, auth, file coordination
- **Asynq (Go)** – Redis-based job queue
- **FFmpeg** – Audio manipulation via `os/exec`
- **bimg/libvips** – Image processing
- **MongoDB** – File metadata & job tracking
- **S3** – File storage
- **CloudFront** – Media delivery (CDN)
- **Asynqmon** – Job dashboard (like BullBoard)

---

### 🧩 Frontend/Backend Processing Boundary

- **Frontend** (WASM/Web APIs):
  - Previews of audio speed, pitch, trim
  - Render waveform in real-time
  - Validate file before upload
- **Backend** (Go):
  - Apply permanent edits: normalization, trimming, format conversion
  - Save final processed version to S3
  - Store original + processed media references

---
