package middleware

import (
	"context"
	"errors"
	"gateway/v1/internal/helper"
	"time"

	redisrate "github.com/go-redis/redis_rate/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type rateLimitMiddleware struct {
	rdb          *redis.Client
	RedisLimiter *redisrate.Limiter
}

func NewRateLimitMiddleware(RedisLimiter *redisrate.Limiter) *rateLimitMiddleware {
	return &rateLimitMiddleware{
		RedisLimiter: RedisLimiter,
	}
}

func (r *rateLimitMiddleware) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		context_, cancel := context.WithTimeout(c.Context(), time.Second*10)
		defer cancel()

		key := c.Get("X-Forwarded-For")
		if key == "" {
			key = c.IP() // fallback
		}

		res, _ := r.RedisLimiter.Allow(context_, key, redisrate.Limit{
			Rate:   3,
			Burst:  3,
			Period: 10 * time.Second,
		})

		if res.Remaining <= 0 {
			return helper.RespondHttpError(c, helper.NewHttpError(fiber.StatusTooManyRequests, errors.New("Request limit exceeded")))
		}
		return c.Next()
	}
}
