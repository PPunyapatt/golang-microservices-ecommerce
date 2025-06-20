package main

import (
	"cart-service/v1/config"
	"cart-service/v1/internal/repository"
	"cart-service/v1/pkg/Database/gorm"
	"cart-service/v1/pkg/Database/postgres"
	"cart-service/v1/pkg/rabbitmq"
	"fmt"
	"log"
)

func main() {
	// s := grpc.NewServer()

	// listener, err := net.Listen("tcp", ":50051")
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }

	cfg, err := config.SetUpEnv()
	if err != nil {
		panic(err)
	}

	// database connection
	// dbGorm, err := gorm.Open(postgres.Open(cfg.Dsn), &gorm.Config{})
	// if err != nil {
	// 	panic(err)
	// }

	// dbSqlx, err := sqlx.Connect("pgx", cfg.Dsn)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	postgresDB, err := postgres.NewPostgresDB(cfg.Dsn)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	gormDB, err := gorm.NewGormConnection(cfg.Dsn)
	if err != nil {
		log.Fatalf("failed to connect to gorm: %v", err)
	}

	// RabbitMQ connection
	rabbitmq, err := rabbitmq.NewRabbitMQConnection("")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}

	cartRepo := repository.NewRepository(postgresDB, gormDB)

	result, err := cartRepo.GetOrCreateCartByUserID("TEST-01")

	if err != nil {
		log.Fatalf("failed to get or create cart: %v", err)
	}
	fmt.Println("Result:", result)

	// items := []*constant.Item{
	// 	{
	// 		ProductID:   1,
	// 		ProductName: "Notebook",
	// 		Quantity:    2,
	// 		Price:       20000,
	// 	},
	// }

	// err = cartRepo.AddItem("TEST-01", items)

	err = cartRepo.RemoveItem("TEST-01", 4, 4)
	if err != nil {
		log.Fatalf("failed to remove item: %v", err)
	}

	fmt.Println("Error:", err)

}
