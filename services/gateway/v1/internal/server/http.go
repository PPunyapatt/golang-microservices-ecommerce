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

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func MapRoutes(ctx context.Context, service *handler.ApiHandler, prometheusMetrics *metrics.Metrics, wg *sync.WaitGroup) {
	app := fiber.New()
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

	// routes
	api.Route(app, service)

	// srv := &http.Server{
	// 	Addr: ":1234",
	// }

	go func() {
		// log.Println("📊 fiber start")
		// if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		// 	log.Printf("❌ fiber error: %v\n", err)
		// }

		// Start the server
		err := app.Listen(":1234")
		if err != nil {
			slog.Error("❌ Fiber error", "error", err.Error())
		}
	}()

	<-ctx.Done()
	// _ = srv.Shutdown(ctx)
	if err := app.Shutdown(); err != nil {
		slog.Error("❌ Fiber shutdown error", "error", err)
	}
	log.Println("🛑 Shutting down fiber...")
	wg.Done()
}
