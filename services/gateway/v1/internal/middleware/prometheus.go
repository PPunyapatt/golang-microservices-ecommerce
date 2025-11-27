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
		status := c.Response().StatusCode()
		statusCode := strconv.Itoa(status)
		method := c.Method()
		path := c.Route().Path
		slog.Info(
			"MiddleWare",
			"method", method,
			"path", path,
			"status", statusCode,
			"duration", duration,
		)
		pm.Http.HttpRequestsTotal.WithLabelValues(statusCode, method, path).Inc()
		pm.Http.HttpRequestDuration.WithLabelValues(method, path, statusCode).Observe(duration)

		if status >= 200 && status < 300 {
			pm.Http.SuccessHttpRequests.Inc()
		} else {
			pm.Http.ErrorHttpRequests.Inc()
		}

		return err
	}
}
