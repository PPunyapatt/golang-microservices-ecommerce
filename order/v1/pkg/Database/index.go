package database

import (
	"log"
	"order/v1/config"

	"github.com/XSAM/otelsql"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Gorm *gorm.DB
	Sqlx *sqlx.DB
}

func InitDatabase(cfg *config.AppConfig) (*Database, error) {
	// ------------- gorm -------------
	dbGorm, err := gorm.Open(postgres.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err.Error())
		return nil, err
	}

	if err := dbGorm.Use(otelgorm.NewPlugin()); err != nil {
		log.Fatal("Failed to use otel: ", err.Error())
		return nil, err
	}

	// ------------- postgres -------------
	dbSqlx, err := otelsql.Open("pgx", cfg.Dsn, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
	))
	if err != nil {
		return nil, err
	}

	if err = otelsql.RegisterDBStatsMetrics(dbSqlx, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
	)); err != nil {
		return nil, err
	}

	log.Println("Connect database success ✅")

	return &Database{
		Gorm: dbGorm,
		Sqlx: sqlx.NewDb(dbSqlx, "pgx"),
	}, nil
}
