package main

import (
	"auth-service/v1/internal/app"
	"auth-service/v1/internal/repository"
	"auth-service/v1/internal/service"
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	database "package/Database"
	"package/config"
	"package/tracer"

	"go.opentelemetry.io/otel"
)

func main() {
	cfg, err := config.SetUpEnv("postgres")
	if err != nil {
		panic(err)
	}

	// ✅ Init tracer
	shutdown := tracer.InitTracer("auth-service")
	defer func() { _ = shutdown(context.Background()) }()

	db, err := database.InitDatabase(cfg)
	if err != nil {
		panic(err)
	}
	authRepo := repository.NewAuthRepository(db.Gorm, db.Sqlx)

	authService := service.NewAuthServer(authRepo, otel.Tracer("inventory-service"), cfg.JwtSecret)

	go func() {
		// เปิด HTTP server ที่ expose /debug/pprof/*
		log.Println("Expose pprof")
		log.Println(http.ListenAndServe("localhost:6061", nil))
	}()

	app.StartgRPCServer(authService)
}
