package repository

import (
	"context"
	"order/v1/internal/constant"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type OrderRepository interface {
	// GetOrder(context.Context, int32) error
	// ListOrder(context.Context, *string) error
	AddOrderItems(context.Context, *gorm.DB, []*constant.OrderItems) error
	AddOrder(context.Context, *gorm.DB, *constant.Order, *int) error
	UpdateStatus(context.Context, int, map[string]interface{}) error
	GetItemsByOrderID(context.Context, int) ([]*constant.InventoryOrder, error)
	CheckAndUpdateStatus(context.Context, int) (bool, error)

	BeginTx() (*gorm.DB, error)

	CreateProduct(context.Context, *constant.Product) error
	UpdateProduct(context.Context, *constant.Product) error
	CalculateTotalPrice(context.Context, map[int]*constant.OrderItems) ([]*constant.OrderItems, error)
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
