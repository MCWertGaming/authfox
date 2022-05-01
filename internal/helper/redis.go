package helper

import (
	"os"

	"github.com/go-redis/redis"
)

func Connect(dbNumber int) *redis.Client {
	// return redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: dbNumber})
	return redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_HOST"), Password: os.Getenv("REDIS_PASS"), DB: dbNumber})
}
