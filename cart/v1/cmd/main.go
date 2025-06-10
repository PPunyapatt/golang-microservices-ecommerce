package main

import (
	"cart-service/v1/config"
	"cart-service/v1/internal/repository"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	dbGorm, err := gorm.Open(postgres.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	dbSqlx, err := sqlx.Connect("pgx", cfg.Dsn)
	if err != nil {
		log.Fatalln(err)
	}

	cartRepo := repository.NewRepository(dbSqlx, dbGorm)

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
