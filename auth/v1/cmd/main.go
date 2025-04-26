package main

import (
	"auth-service/v1/config"
	"auth-service/v1/internal/api"
	"auth-service/v1/internal/repository"
	"auth-service/v1/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// server
	app := fiber.New()
	c := cors.New(cors.Config{
		AllowOrigins: "http://localhost:3030",
		AllowHeaders: "Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	})

	app.Use(c)

	cfg, err := config.SetUpEnv()
	if err != nil {
		panic(err)
	}

	// database connection
	dbGorm, err := gorm.Open(postgres.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	dbSqlx, err := sqlx.Connect("pgx", cfg.Dsn)
	if err != nil {
		log.Fatalln(err)
	}

	authRepo := repository.NewAuthRepository(dbGorm, dbSqlx)
	authService := service.NewAuthService(authRepo)

	// routes
	api.Route(app, authService)

	// go func() {
	err = app.Listen(":" + cfg.ServerPort)
	if err != nil {
		log.Fatal(err)
	}
}
