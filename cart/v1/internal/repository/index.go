package repository

import (
	"cart/v1/internal/constant"
	"cart/v1/proto/cart"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type Repository struct {
	PostgresDB *sqlx.DB
	GormDB     *gorm.DB
}

type CartRepository interface {
	CreateCart(userID string) error
	GetOrCreateCartByUserID(userID string) (*constant.Cart, error)
	AddItem(userID string, items []*constant.Item) error
	RemoveItem(userID string, cartID, itemID int) error
	RemoveCart(userID string) error
	GetItemsByUserID(userID string, pagination *constant.Pagination) ([]*cart.CartItem, error)
}

func NewRepository(postgresDB *sqlx.DB, gormDB *gorm.DB) CartRepository {
	return &Repository{
		PostgresDB: postgresDB,
		GormDB:     gormDB,
	}
}
