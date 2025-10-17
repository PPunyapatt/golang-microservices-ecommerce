package main

import (
	"context"
	"gateway/v1/internal/api"
	"gateway/v1/internal/api/handler"
	"gateway/v1/internal/helper"
	"log"
	"net/http"
	_ "net/http/pprof"
	"package/tracer"

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

	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, skipping...")
	}

	shutdown := tracer.InitTracer("gateway")
	defer func() { _ = shutdown(context.Background()) }()

	// Add logger middleware
	app.Use(logger.New())
	app.Use(c)
	app.Use(otelfiber.Middleware())

	conn := helper.NewClientsGRPC()

	log.Println("Connected to all gRPC server")

	// routes
	service := handler.ServiceNew(conn)
	api.Route(app, service)

	go func() {
		// เปิด HTTP server ที่ expose /debug/pprof/*
		http.ListenAndServe(":6060", nil)
	}()

	// Start the server
	err := app.Listen(":1234")
	if err != nil {
		log.Fatal(err)
	}
}
