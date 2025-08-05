package main

import (
	"context"
	"gateway/v1/internal/api"
	"gateway/v1/internal/api/handler"
	"gateway/v1/internal/helper"
	"log"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()
	c := cors.New(cors.Config{
		AllowOrigins: "http://localhost:3030",
		AllowHeaders: "Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	})

	shutdown := helper.InitTracer("gateway")
	defer func() { _ = shutdown(context.Background()) }()

	// Add logger middleware
	app.Use(logger.New())
	app.Use(c)
	app.Use(otelfiber.Middleware())

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err.Error())
	}

	conn := helper.NewClientsGRPC()

	log.Println("Connected to all gRPC server")

	// routes
	service := handler.ServiceNew(conn)
	api.Route(app, service)

	// Start the server
	err = app.Listen(":1234")
	if err != nil {
		log.Fatal(err)
	}
}
