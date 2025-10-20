package middleware

import (
	"log/slog"
	"strconv"
	"time"

	"package/metrics"

	"github.com/gofiber/fiber/v2"
)

func PrometheusMiddleware(pm *metrics.Metrics) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		duration := time.Since(start).Seconds()
		statusCode := strconv.Itoa(c.Response().StatusCode())
		method := c.Method()
		path := c.Route().Path

		pm.HttpRequestsTotal.WithLabelValues(method, path, statusCode).Inc()
		pm.HttpRequestDuration.WithLabelValues(method, path, statusCode).Observe(duration)

		slog.Info(
			"method", method,
			slog.Float64("response time", duration),
			"Path", path,
			"Status", statusCode,
		)
		return err
	}
}
