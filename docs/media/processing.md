# Media Processing

TuneSlap processes uploaded media files to normalize formats, optimize file sizes, and extract metadata. Processing happens in background tasks after upload.

## Processing Pipeline

When a media file is uploaded:

1. Raw file is stored in the User Uploads bucket
2. A background task is queued (see [Async Operations](./async-operations.md))
3. The file is downloaded to a temp directory
4. Processing is applied based on media type
5. Processed file is uploaded to the Media bucket
6. Media record is updated with new metadata

## Audio Processing

Audio files are processed using FFmpeg via the `mempeg/audio` and `mempeg/ffmpeg` packages.

### Default Normalization

All audio files are normalized with these settings:

| Setting | Value | Description |
|---------|-------|-------------|
| Loudness | -16 LUFS | EBU R128 loudness standard |
| True Peak | -1.5 dBTP | Prevents clipping |
| Loudness Range | 11 LU | Dynamic range target |
| Sample Rate | 44100 Hz | CD quality |
| Silence Removal | -50 dB threshold | Trims leading silence |

### Output Formats

Audio is encoded to:
- **WebM (Opus)** at 128 kbps - Primary format for web playback
- **MP3** at 128 kbps - Fallback for older browsers

### AudioProcessingParams

The `AudioProcessingParams` struct controls audio processing:

```go
type AudioProcessingParams struct {
    ContentType   string   // Output MIME type: "audio/webm", "audio/mp3"
    TrimStart     float64  // Trim from start (seconds)
    TrimEnd       float64  // Trim from end (seconds)
    Normalize     bool     // Apply loudness normalization
    FadeIn        float64  // Fade in duration (seconds)
    FadeOut       float64  // Fade out duration (seconds)
    Speed         float64  // Playback speed multiplier (e.g., 1.25)
    Pitch         float64  // Pitch shift in semitones
    OutputFormats []string // e.g., ["webm", "mp3"]
}
```

**Example request:**

```json
{
  "processingParams": {
    "audio": {
      "contentType": "audio/webm",
      "normalize": true,
      "trimStart": 0.5,
      "fadeIn": 0.2,
      "fadeOut": 0.5
    }
  }
}
```

### Extracted Metadata

After processing, these fields are populated on the media record:

- `duration` - Length in seconds
- `fileSize` - Processed file size in bytes
- `contentType` - MIME type of processed file
- `processedUrl` - URL to the processed file

## Image Processing

Images are processed using [bimg](https://github.com/h2non/bimg) (libvips bindings) via the `mempeg/image` package.

### Default Normalization

All images are normalized with these settings:

| Setting | Value | Description |
|---------|-------|-------------|
| Max Dimensions | 500x500 | Resized to fit |
| Quality | 90 | JPEG/WebP quality |
| Format | WebP | Output format |
| Compression | 6 | WebP compression level |
| Strip Metadata | Yes | Removes EXIF data |

### ImageProcessingParams

The `ImageProcessingParams` struct controls image processing:

```go
type ImageProcessingParams struct {
    ResizeTo     [2]int // [width, height]
    Format       int    // Output format (WebP, PNG, etc.)
    Crop         [4]int // [x, y, width, height]
    AspectRatio  string // "16:9", "1:1", etc.
    ApplyFilters string // "grayscale", "blur", etc.
}
```

**Example request:**

```json
{
  "processingParams": {
    "image": {
      "resizeTo": [800, 600],
      "applyFilters": "grayscale"
    }
  }
}
```

### Extracted Metadata

After processing, these fields are populated on the media record:

- `dimensions` - [width, height] in pixels
- `fileSize` - Processed file size in bytes
- `contentType` - MIME type of processed file
- `processedUrl` - URL to the processed file

## Processing Status

Media records have a `status` field that tracks processing state:

| Status | Description |
|--------|-------------|
| `pending` | Uploaded, waiting for processing |
| `processing` | Currently being processed |
| `done` | Processing complete |
| `error` | Processing failed |

The `processingActivity` array provides a history of status changes with timestamps and messages.

## Dependencies

### FFmpeg

Required for audio processing. Must be installed on the server.

```bash
# Ubuntu/Debian
apt-get install ffmpeg

# macOS
brew install ffmpeg

# Docker (included in server image)
```

### libvips

Required for image processing (via bimg).

```bash
# Ubuntu/Debian
apt-get install libvips-dev

# macOS
brew install vips

# Docker (included in server image)
```

## Related

- [Async Operations](./async-operations.md) - How processing tasks are queued
- [Storage](./storage.md) - Where processed files are stored
- [Upload Flow](./upload-flow.md) - How files get to processing
