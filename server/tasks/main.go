package tasks

import (
	"fmt"
	"log"
	"os"
	"tuneslap/config"

	"github.com/hibiken/asynq"
)

var Client *asynq.Client
var scheduler *asynq.Scheduler

func InitClient() {
	fmt.Println("Starting Asynq Client...")
	Client = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}

func GetClient() *asynq.Client {
	return Client
}

func CloseClient() {
	if Client != nil {
		Client.Close()
	}
}

// StartWorker starts the asynq worker server to process background tasks
func StartWorker() {
	log.Println("[Worker] Starting Asynq Worker...")
	log.Printf("[Worker] Connecting to Redis at: %s", os.Getenv("REDIS_URL"))

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     os.Getenv("REDIS_URL"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"media":   10,
				"default": 5,
			},
			Logger: &workerLogger{},
		},
	)

	mux := asynq.NewServeMux()
	RegisterMediaProcessTasks(mux)
	RegisterCleanupTasks(mux)

	log.Println("[Worker] Worker server starting, waiting for tasks...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("[Worker] Could not run asynq server: %v", err)
	}
}

// workerLogger implements asynq.Logger for better visibility
type workerLogger struct{}

func (l *workerLogger) Debug(args ...interface{}) {
	log.Println("[Worker DEBUG]", fmt.Sprint(args...))
}

func (l *workerLogger) Info(args ...interface{}) {
	log.Println("[Worker INFO]", fmt.Sprint(args...))
}

func (l *workerLogger) Warn(args ...interface{}) {
	log.Println("[Worker WARN]", fmt.Sprint(args...))
}

func (l *workerLogger) Error(args ...interface{}) {
	log.Println("[Worker ERROR]", fmt.Sprint(args...))
}

func (l *workerLogger) Fatal(args ...interface{}) {
	log.Fatal("[Worker FATAL]", fmt.Sprint(args...))
}

// StartScheduler starts the asynq scheduler for periodic tasks (only in demo mode)
func StartScheduler() {
	if !config.IsDemoMode() {
		log.Println("[Scheduler] Not in demo mode, skipping scheduler")
		return
	}

	log.Println("[Scheduler] Starting Asynq Scheduler for demo mode...")

	scheduler = asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:     os.Getenv("REDIS_URL"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		},
		&asynq.SchedulerOpts{
			Logger: &workerLogger{},
		},
	)

	// Schedule demo cleanup every hour
	task, err := NewDemoCleanupTask()
	if err != nil {
		log.Printf("[Scheduler] Failed to create cleanup task: %v", err)
		return
	}

	entryID, err := scheduler.Register("@every 1h", task)
	if err != nil {
		log.Printf("[Scheduler] Failed to register cleanup task: %v", err)
		return
	}

	log.Printf("[Scheduler] Registered demo cleanup task with entry ID: %s", entryID)

	if err := scheduler.Run(); err != nil {
		log.Printf("[Scheduler] Could not run scheduler: %v", err)
	}
}

// CloseScheduler stops the scheduler
func CloseScheduler() {
	if scheduler != nil {
		scheduler.Shutdown()
	}
}
