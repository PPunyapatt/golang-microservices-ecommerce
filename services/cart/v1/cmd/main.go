package main

import (
	"cart/v1/internal/app"
	"cart/v1/internal/repository"
	"cart/v1/internal/service"
	"config-service"
	"context"
	database "package/Database"
	"package/rabbitmq"
	"package/tracer"
)

func main() {
	cfg, err := config.SetUpEnv()
	if err != nil {
		panic(err)
	}

	// ✅ Init tracer
	shutdown := tracer.InitTracer("order-service")
	defer func() { _ = shutdown(context.Background()) }()

	// database connection
	db, err := database.InitDatabase(cfg)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.Gorm.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()
	defer db.Sqlx.Close()

	// RabbitMQ Connection
	conn, err := rabbitmq.NewRabbitMQConnection(cfg.RabbitMQUrl)
	if err != nil {
		panic(err)
	}

	cartRepo := repository.NewRepository(db.Sqlx, db.Gorm, db.Mongo)

	cartService, cartServerRPC := service.NewCartServer(cartRepo)

	app.InitConsumer(cartService, conn)

	app.StartgRPCServer(cartServerRPC)
}
