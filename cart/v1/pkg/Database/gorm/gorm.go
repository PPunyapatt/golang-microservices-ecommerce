package gorm

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	_defaultConnAttempts = 3
	_defaultConnTimeout  = time.Second
)

type gormDB struct {
	connAttempts int
	connTimeout  time.Duration

	db *gorm.DB
}

func NewGormConnection(dsn string) (*gorm.DB, error) {
	gorm_ := &gormDB{
		connAttempts: _defaultConnAttempts,
		connTimeout:  2 * _defaultConnTimeout,
	}

	var err error
	for gorm_.connAttempts > 0 {
		gorm_.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("attempt %d: failed to open DB: %v", gorm_.connAttempts, err)
		} else {
			log.Printf("📰 connected to gorm 🎉")
			break
		}
	}

	if err != nil {
		return nil, err
	}

	return gorm_.db, nil
}
