package postgres

import (
	"cart-service/v1/config"
	"log"

	"github.com/jmoiron/sqlx"
)

func NewConnection(config *config.AppConfig) (*sqlx.DB, error) {
	dbSqlx, err := sqlx.Connect("pgx", config.Dsn)
	if err != nil {
		log.Fatalln(err)
	}
	return dbSqlx, nil
}
