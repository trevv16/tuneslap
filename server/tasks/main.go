package tasks

import (
	"fmt"
	"log"
	"os"

	"github.com/hibiken/asynq"
)

var Client *asynq.Client

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
