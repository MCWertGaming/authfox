package helper

import (
	"github.com/go-redis/redis"
)

func Connect(dbNumber int) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: dbNumber})
}
