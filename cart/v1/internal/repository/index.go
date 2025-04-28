package repository

import (
	"cart-service/v1/config"
	gormdb "cart-service/v1/internal/repository/gorm"
	"cart-service/v1/internal/repository/postgres"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type Repository struct {
	PostgresDB *sqlx.DB
	GormDB     *gorm.DB
}

func SetupDatabase(config *config.AppConfig) (*sqlx.DB, *gorm.DB, error) {
	// Setup PostgresDB connection
	postgresDB, err := postgres.NewConnection(config)
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gormdb.NewConnection(config)
	if err != nil {
		return nil, nil, err
	}

	return postgresDB, gormDB, nil
}

func NewRepository(postgresDB *sqlx.DB, gormDB *gorm.DB) cartRepository {
	return &Repository{
		PostgresDB: postgresDB,
		GormDB:     gormDB,
	}
}
