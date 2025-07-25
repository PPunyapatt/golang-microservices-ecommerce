package main

import (
	"gateway/v1/internal/api"
	"gateway/v1/internal/api/handler"
	"gateway/v1/proto/auth"
	"gateway/v1/proto/cart"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	app := fiber.New()
	c := cors.New(cors.Config{
		AllowOrigins: "http://localhost:3030",
		AllowHeaders: "Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	})

	// Add logger middleware
	app.Use(logger.New())

	app.Use(c)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err.Error())
	}

	cc, err := grpc.NewClient("localhost:1024", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer cc.Close()

	cartClient := cart.NewCartServiceClient(cc)
	authClient := auth.NewAuthServiceClient(cc)

	log.Println("Connected to gRPC server")

	// routes
	service := handler.ServiceNew(authClient, cartClient)
	api.Route(app, service)

	// Start the server
	err = app.Listen(":1234")
	if err != nil {
		log.Fatal(err)
	}
}
