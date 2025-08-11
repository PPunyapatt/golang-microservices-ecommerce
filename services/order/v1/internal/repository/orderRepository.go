package repository

import (
	"context"

	"gorm.io/gorm"
)

// func (repo *orderRepository) GetOrder(ctx context.Context, orderID int32) error {
// 	return nil
// }

// func (repo *orderRepository) ListOrder(ctx context.Context, status *string) error {
// 	return nil
// }

func (repo *orderRepository) BeginTx() (*gorm.DB, error) {
	tx := repo.gorm.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return tx, nil
}

func (repo *orderRepository) AddOrderItems(tx *gorm.DB, ctx context.Context) error {
	return nil
}

func (repo *orderRepository) AddOrder(tx *gorm.DB, ctx context.Context) (int, error) {
	return nil
}

func (repo *orderRepository) ChangeStatus(ctx context.Context, orderStatus string, paymentStatus string) error {
	return nil
}
