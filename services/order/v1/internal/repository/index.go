package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type OrderRepository interface {
	// GetOrder(context.Context, int32) error
	// ListOrder(context.Context, *string) error
	AddOrderItems(context.Context) error
	AddOrder(context.Context) error
	ChangeStatus(context.Context, string, string) error
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
