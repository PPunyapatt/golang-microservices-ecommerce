package repository

import (
	"cart/v1/internal/constant"
	"cart/v1/proto/cart"
	"context"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
)

type Repository struct {
	PostgresDB *sqlx.DB
	GormDB     *gorm.DB
	MongoDB    *mongo.Client
}

type CartRepository interface {
	CreateCart(userID string) error
	GetOrCreateCartByUserID(ctx context.Context, userID string, pagination *constant.Pagination) ([]*cart.CartItem, error)
	AddItem(userID string, items []*constant.Item) error
	RemoveItem(userID string, cartID, itemID int) error
	RemoveCart(userID string) error
	GetItemsByUserID(userID string, pagination *constant.Pagination) ([]*cart.CartItem, error)
	DeleteCart(context.Context, string) error
}

func NewRepository(postgresDB *sqlx.DB, gormDB *gorm.DB, mongoDB *mongo.Client) CartRepository {
	return &Repository{
		PostgresDB: postgresDB,
		GormDB:     gormDB,
		MongoDB:    mongoDB,
	}
}
