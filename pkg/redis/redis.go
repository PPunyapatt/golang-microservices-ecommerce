package redis

import (
	"context"
	"fmt"
	"log/slog"

	"package/config"

	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.AppConfig) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		slog.Error("Connect to redis failed",
			"err", fmt.Errorf("failed to connect to redis: %w", err),
		)
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return rdb, nil
}
