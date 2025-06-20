package postgres

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	_defaultConnAttempts = 3
	_defaultConnTimeout  = time.Second
)

type DBConnString string

type postgres struct {
	connAttempts int
	connTimeout  time.Duration

	db *sqlx.DB
}

func NewPostgresDB(dsn string) (*sqlx.DB, error) {
	// slog.Info("CONN", "connect string", url)

	pg := &postgres{
		connAttempts: _defaultConnAttempts,
		connTimeout:  2 * _defaultConnTimeout,
	}

	var err error
	for pg.connAttempts > 0 {
		pg.db, err = sqlx.Connect("pgx", dsn)
		if err != nil {
			log.Printf("attempt %d: failed to open DB: %v", pg.connAttempts, err)
		} else {
			log.Printf("📰 connected to postgresdb 🎉")
			break
		}

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, err
	}

	// slog.Info("📰 connected to postgresdb 🎉")

	return pg.db, nil
}
