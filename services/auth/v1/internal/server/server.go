package server

import (
	"auth-service/v1/internal/app"
	"auth-service/v1/internal/repository"
	"auth-service/v1/internal/service"
	"context"
	"log/slog"
	"os"
	"os/signal"
	database "package/Database"
	"package/config"
	"package/metrics"
	"package/tracer"
	"sync"
	"syscall"

	"go.opentelemetry.io/otel"
)

type server struct {
	cfg *config.AppConfig
}

func NewServer(cfg *config.AppConfig) *server {
	return &server{cfg}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// ✅ Init tracer
	shutdown := tracer.InitTracer("auth-service")
	defer func() { _ = shutdown(ctx) }()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// Initialize middleware
	prometheusMetrics := metrics.NewMetrics()

	db, err := database.InitDatabase(s.cfg)
	if err != nil {
		return err
	}
	authRepo := repository.NewAuthRepository(db.Gorm, db.Sqlx)

	authService := service.NewAuthServer(authRepo, otel.Tracer("auth-service"), s.cfg.JwtSecret)

	var wg sync.WaitGroup
	wg.Add(2)
	go app.StartgRPCServer(ctx, &wg, authService)
	// Start Prometheus endpoint
	go prometheusMetrics.PrometheusHttp(ctx, &wg)
	wg.Wait()

	return nil
}
