package gorm

import (
	"cart-service/v1/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(config *config.AppConfig) (*gorm.DB, error) {
	dbGorm, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return dbGorm, nil
}
