# Async Operations

TuneSlap uses [Asynq](https://github.com/hibiken/asynq) for background task processing, backed by Redis. This allows media processing to happen asynchronously without blocking API requests.

## Architecture

```
┌──────────┐     enqueue      ┌──────────┐     dequeue      ┌──────────┐
│ API      │ ────────────────▶│  Redis   │◀────────────────▶│ Worker   │
│ Handler  │                  │  Queue   │                  │ Process  │
└──────────┘                  └──────────┘                  └──────────┘
```

1. **API Handler** - Receives upload confirmation and enqueues a processing task
2. **Redis Queue** - Stores pending tasks with payload data
3. **Worker Process** - Picks up tasks and runs the processing logic

## Task Types

### `media:process`

The primary task type for processing uploaded media files.

**Payload:**

```go
type MediaProcessPayload struct {
    MediaID          primitive.ObjectID      // ID of the media record
    UserID           primitive.ObjectID      // Owner of the media
    ProcessingParams models.ProcessingParams // Audio or image params
}
```

**Lifecycle:**

1. Task is created when media record is created
2. Worker picks up the task from Redis
3. Media status is set to `processing`
4. File is downloaded, processed, and re-uploaded
5. Media status is set to `done` (or `error` on failure)

## Enqueueing Tasks

Tasks are enqueued from API handlers using the Asynq client:

```go
import "tuneslap/tasks"

// Create the task
task, err := tasks.NewMediaProcessTask(mediaId, userId, processingParams)
if err != nil {
    return err
}

// Enqueue it
_, err = tasks.GetClient().Enqueue(task)
if err != nil {
    return err
}
```

## Task Handler

The task handler in `server/tasks/mediaTasks.go` processes each task:

```go
func HandleMediaProcessTask(ctx context.Context, task *asynq.Task) error {
    // 1. Unmarshal payload
    var payload MediaProcessPayload
    json.Unmarshal(task.Payload(), &payload)

    // 2. Get media record
    media, _ := mediaRepo.GetByIdUnscoped(payload.MediaID)

    // 3. Update status to "processing"
    mediaRepo.UpdateMediaUnscoped(payload.MediaID, &models.Media{
        Status: models.ProcessingStatusProcessing,
    })

    // 4. Process based on media type
    switch media.MediaType {
    case "audio":
        processedAudio, _ := audio.ProcessAudio(media, *payload.ProcessingParams.Audio)
        // Update with processed data
    case "image":
        processedImage, _ := image.ProcessImage(media, *payload.ProcessingParams.Image)
        // Update with processed data
    }

    // 5. Update status to "done"
    mediaRepo.UpdateMediaUnscoped(payload.MediaID, &updateData)

    return nil
}
```

## Redis Configuration

The Asynq client connects to Redis using environment variables:

| Variable | Required | Description |
|----------|----------|-------------|
| `REDIS_URL` | Yes | Redis address (e.g., `redis:6379` or `localhost:6379`) |
| `REDIS_PASSWORD` | No | Redis password if authentication is enabled |

## Processing Activity

Each media record maintains a `processingActivity` array that logs status changes:

```json
{
  "processingActivity": [
    {
      "status": "pending",
      "message": "Upload received",
      "createdAt": "2024-01-15T10:00:00Z"
    },
    {
      "status": "processing",
      "message": "Processing started",
      "createdAt": "2024-01-15T10:00:05Z"
    },
    {
      "status": "done",
      "message": "Processing completed",
      "createdAt": "2024-01-15T10:00:15Z"
    }
  ]
}
```

## Error Handling

If processing fails:

1. The task handler returns an error
2. Asynq automatically retries the task (configurable)
3. After max retries, the task is moved to the dead queue
4. Media status remains `processing` or is set to `error`

To inspect failed tasks, use the Asynq CLI or web UI.

## Monitoring

### Asynq CLI

Install the Asynq CLI for queue inspection:

```bash
go install github.com/hibiken/asynq/tools/asynq@latest

# List queues
asynq queue ls

# List pending tasks
asynq task ls --queue=default --state=pending

# List failed tasks
asynq task ls --queue=default --state=archived
```

### Redis CLI

You can also inspect the queue directly in Redis:

```bash
redis-cli

# List all keys
KEYS asynq:*

# Check queue length
LLEN asynq:default:pending
```

## Development

In development, the worker runs as part of the main server process. The server initializes the Asynq client on startup:

```go
// server/tasks/main.go
func InitClient() {
    Client = asynq.NewClient(asynq.RedisClientOpt{
        Addr:     os.Getenv("REDIS_URL"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
    })

    RegisterMediaProcessTasks()
}
```

## Related

- [Processing](./processing.md) - What happens inside the task
- [Upload Flow](./upload-flow.md) - How tasks get queued
- [Storage](./storage.md) - Where processed files are stored
