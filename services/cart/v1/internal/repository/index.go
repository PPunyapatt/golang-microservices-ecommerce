package repository

import (
	"cart/v1/internal/constant"
	"cart/v1/proto/cart"
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	// PostgresDB *sqlx.DB
	// GormDB     *gorm.DB
	MongoDB *mongo.Client
	Logger  *slog.Logger
}

type CartRepository interface {
	GetOrCreateCartByUserID(ctx context.Context, userID string, pagination *constant.Pagination) ([]*cart.StoreItems, error)
	AddItem(ctx context.Context, userID string, items []*constant.StoreItems) error
	RemoveItem(ctx context.Context, userID string, itemID int) error
	RemoveCart(ctx context.Context, userID string) error
}

func NewRepository(mongoDB *mongo.Client, logger *slog.Logger) CartRepository {
	return &Repository{
		// PostgresDB: postgresDB,
		// GormDB:     gormDB,
		MongoDB: mongoDB,
		Logger:  logger,
	}
}
