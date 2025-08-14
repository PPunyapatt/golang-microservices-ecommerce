package repository

import (
	"context"
	"log"
	"order/v1/internal/constant"

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

func (repo *orderRepository) AddOrder(tx *gorm.DB, ctx context.Context) (*int, error) {
	return nil, nil
}

func (repo *orderRepository) ChangeStatus(ctx context.Context, orderStatus string, paymentStatus string) error {
	return nil
}

func (repo *orderRepository) CreateProduct(ctx context.Context, product *constant.Product) error {
	result := repo.gorm.Omit("updated_at").Create(product)
	if result.Error != nil {
		log.Println("Failed")
		return result.Error
	}

	log.Println("Success")

	return nil
}

func (repo *orderRepository) UpdateProduct(ctx context.Context, product *constant.Product) error {
	updateData := map[string]interface{}{}
	if product.StoreID != nil {
		updateData["store_id"] = product.StoreID
	}

	if product.ProductName != nil {
		updateData["product_name"] = product.ProductName
	}

	if product.Price != nil {
		updateData["price"] = product.Price
	}

	result := repo.gorm.Model(&constant.Product{}).Where("product_id = ?", product.ProductID).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
