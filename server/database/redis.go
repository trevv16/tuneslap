package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	fmt.Println("Starting Redis...")
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}

func GetRedisClient() *redis.Client {
	return RedisClient
}

func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
	}
}

func SetCache(key string, value interface{}, ttl time.Duration) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client is not initialized")
	}
	return RedisClient.Set(context.Background(), key, value, ttl).Err()
}

func GetCache(key string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("Redis client is not initialized")
	}
	return RedisClient.Get(context.Background(), key).Result()
}

func DeleteCache(key string) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client is not initialized")
	}
	return RedisClient.Del(context.Background(), key).Err()
}
