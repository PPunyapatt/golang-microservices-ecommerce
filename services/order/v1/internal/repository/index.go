package repository

import (
	"context"
	"order/v1/internal/constant"
	"order/v1/proto/order"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type OrderRepository interface {
	// GetOrder(context.Context, int32) error
	ListOrder(ctx context.Context, req *constant.ListOrderRequest, pagination *constant.Pagination) ([]*order.Orders, error)
	AddOrderItems(context.Context, *gorm.DB, []*constant.OrderItems) error
	AddOrder(context.Context, *gorm.DB, *constant.Order, *int) error
	UpdateStatus(context.Context, int, map[string]interface{}) error
	GetItemsByOrderID(context.Context, int) ([]*constant.InventoryOrder, error)
	CheckAndUpdateStatus(context.Context, int) error

	BeginTx() (*gorm.DB, error)

	CreateProduct(context.Context, *constant.Product) error
	UpdateProduct(context.Context, *constant.Product) error
	CalculateTotalPrice(context.Context, map[int]*constant.OrderItems) ([]*constant.OrderItems, error)
	CheckOrderStatus(context.Context, int) (bool, error)
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
