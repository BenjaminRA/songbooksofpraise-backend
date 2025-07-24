package redisdb

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v7"
)

var client *redis.Client

func GetRedisConnection() *redis.Client {
	if client == nil {
		host := os.Getenv("REDIS_HOST")
		port := os.Getenv("REDIS_PORT")
		password := os.Getenv("REDIS_PASSWORD")

		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: password,
		})

		_, err := client.Ping().Result()
		if err != nil {
			panic(err)
		}
	}

	return client
}
