package database

import (
	"log"
	"package/config"

	"github.com/XSAM/otelsql"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Gorm  *gorm.DB
	Sqlx  *sqlx.DB
	Mongo *mongo.Client
}

func InitDatabase(cfg *config.AppConfig) (*Database, error) {
	database := &Database{}
	// -------------------- gorm --------------------
	if cfg.Dsn != "" {
		dbGorm, err := gorm.Open(postgres.Open(cfg.Dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to database: ", err.Error())
			return nil, err
		}

		if err := dbGorm.Use(otelgorm.NewPlugin()); err != nil {
			log.Fatal("Failed to use otel: ", err.Error())
			return nil, err
		}
		database.Gorm = dbGorm

		// -------------------- postgres --------------------
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

		database.Sqlx = sqlx.NewDb(dbSqlx, "pgx")
	}

	// -------------------- MongoDB --------------------
	if cfg.MongoURL != "" {
		opts := options.Client()
		opts.ApplyURI(cfg.MongoURL)
		mongoClient, err := mongo.Connect(opts)
		if err != nil {
			return nil, err
		}
		database.Mongo = mongoClient
	}

	log.Println("Connect database success ✅")

	return database, nil
}
