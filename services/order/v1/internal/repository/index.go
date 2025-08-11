package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type OrderRepository interface {
	// GetOrder(context.Context, int32) error
	// ListOrder(context.Context, *string) error
	AddOrderItems(*gorm.DB, context.Context) error
	AddOrder(*gorm.DB, context.Context) (int, error)
	ChangeStatus(context.Context, string, string) error
	BeginTx() (*gorm.DB, error)
}

type orderRepository struct {
	gorm *gorm.DB
	sqlx *sqlx.DB
}

func NewOrderRepository(gorm *gorm.DB, sqlx *sqlx.DB) OrderRepository {
	return &orderRepository{
		gorm: gorm,
		sqlx: sqlx,
	}
}
