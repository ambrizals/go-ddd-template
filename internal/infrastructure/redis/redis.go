package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/user/go-ddd-template/internal/config"
)

var RedisClient *redis.Client

func InitRedis(cfg *config.Config) error {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	RedisClient = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return RedisClient.Ping(context.Background()).Err()
}
