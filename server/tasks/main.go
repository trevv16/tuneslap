package tasks

import (
	"fmt"
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

	// register tasks
	RegisterMediaProcessTasks()
}

func GetClient() *asynq.Client {
	return Client
}

func CloseClient() {
	if Client != nil {
		Client.Close()
	}
}
