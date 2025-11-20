package server

import (
	"context"
	"gateway/v1/internal/api"
	"gateway/v1/internal/api/handler"
	"gateway/v1/internal/middleware"
	"log"
	"log/slog"
	"sync"

	"package/metrics"

	redisrate "github.com/go-redis/redis_rate/v10"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/redis/go-redis/v9"
)

func MapRoutes(ctx context.Context, service *handler.ApiHandler, prometheusMetrics *metrics.Metrics, wg *sync.WaitGroup, rdb *redis.Client) {
	app := fiber.New()
	errCh := make(chan error, 1)
	c := cors.New(cors.Config{
		AllowOrigins: "http://localhost:3030",
		AllowHeaders: "Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	})

	// Add logger middleware
	app.Use(c)
	app.Use(otelfiber.Middleware())
	app.Use(middleware.PrometheusMiddleware(prometheusMetrics))
	app.Use(logger.New())
	app.Use(middleware.NewRateLimitMiddleware(redisrate.NewLimiter(rdb)).Handler())

	// routes
	api.Route(app, service)

	go func() {
		// Start the server
		err := app.Listen(":1234")
		if err != nil {
			slog.Error("❌ Fiber error", "error", err.Error())
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("🛑 context done, shutting down fiber")
	case err := <-errCh:
		slog.Error("❌ Fiber error", "error", err.Error())
	}

	<-ctx.Done()
	if err := app.Shutdown(); err != nil {
		slog.Error("❌ Fiber shutdown error", "error", err)
	}
	log.Println("🛑 Shutting down fiber...")
	wg.Done()
}
