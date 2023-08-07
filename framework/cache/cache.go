package cache

import (
	"github.com/redis/go-redis/v9"
)

func InitCache(port string) *redis.Client {
	return redis.NewClient(
		&redis.Options{
			Addr:     "localhost:" + port,
			Password: "",
			DB:       0,
		},
	)
}
