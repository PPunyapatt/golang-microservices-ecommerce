package main

import (
	"inventories/v1/config"
	"inventories/v1/internal/repository"
	"inventories/v1/internal/services"
	"inventories/v1/proto/Inventory"
	"log"
	"net"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.SetUpEnv()
	if err != nil {
		panic(err)
	}

	// database connection
	dbGorm, err := gorm.Open(postgres.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err.Error())
	}

	dbSqlx, err := sqlx.Connect("pgx", cfg.Dsn)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err.Error())
	}

	log.Println("Connect database success")

	inventoryRepo := repository.NewInventoryRepository(dbGorm, dbSqlx)

	s := grpc.NewServer()

	// ✅ Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	listener, err := net.Listen("tcp", ":1026")
	if err != nil {
		panic(err)
	}

	Inventory.RegisterInventoryServiceServer(s, services.NewInventoryServer(inventoryRepo))

	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}
}
