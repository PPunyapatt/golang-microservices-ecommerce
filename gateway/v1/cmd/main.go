package main

import (
	"gateway/v1/internal/api"
	"gateway/v1/internal/api/handler"
	"gateway/v1/proto/auth"
	"gateway/v1/proto/cart"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Use(c)

	cc, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	err = app.Listen(":1024")
	if err != nil {
		log.Fatal(err)
	}
}
