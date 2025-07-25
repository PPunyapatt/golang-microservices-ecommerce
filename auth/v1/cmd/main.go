package main

import (
	"auth-service/v1/config"
	"auth-service/v1/internal/repository"
	"auth-service/v1/internal/service"
	"auth-service/v1/proto/auth"
	"log"
	"net"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
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

	authRepo := repository.NewAuthRepository(dbGorm, dbSqlx)

	s := grpc.NewServer()

	listener, err := net.Listen("tcp", ":1024")
	if err != nil {
		panic(err)
	}

	auth.RegisterAuthServiceServer(s, service.NewAuthServer(authRepo))

	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}

	log.Println("Auth service is running on port 1024")
}
