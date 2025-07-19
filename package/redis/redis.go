package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // ubah jika pakai env
		Password: "",               // no password by default
		DB:       0,                // default DB
	})
	return rdb
}

func GetContext() context.Context {
	return ctx
}
