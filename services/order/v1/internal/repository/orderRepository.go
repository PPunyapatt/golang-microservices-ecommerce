package repository

import (
	"context"
	"fmt"
	"log"
	"order/v1/internal/constant"
	"strings"

	"github.com/pkg/errors"
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

func (repo *orderRepository) AddOrderItems(ctx context.Context, tx *gorm.DB, items []*constant.OrderItems) error {
	result := tx.WithContext(ctx).Create(items)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *orderRepository) AddOrder(ctx context.Context, tx *gorm.DB, order *constant.Order) (*int, error) {
	result := tx.WithContext(ctx).Create(order)
	if result.Error != nil {
		return nil, result.Error
	}

	return &order.OrderID, nil
}

func (repo *orderRepository) UpdateStatus(ctx context.Context, orderID int, args ...string) error {
	result := repo.gorm.Model(&constant.Order{}).WithContext(ctx).Where("id = ?", orderID).Update(args[0], args[1])
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *orderRepository) CreateProduct(ctx context.Context, product *constant.Product) error {
	result := repo.gorm.WithContext(ctx).Omit("updated_at").Create(product)
	if result.Error != nil {
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

	result := repo.gorm.Model(&constant.Product{}).WithContext(ctx).Where("product_id = ?", product.ProductID).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *orderRepository) CalculateTotalPrice(ctx context.Context, items map[int]*constant.OrderItems) ([]*constant.OrderItems, error) {
	query := `
		SELECT
			p.product_id,
			p.product_name,
			p.price as unit_price,
			(p.price * i.quantity) AS total_price
		FROM order_products p
		INNER JOIN (
			VALUES
				%s
		) AS i(product_id, quantity)
			ON i.product_id = p.product_id
	`

	args := []interface{}{}
	valueStrings := make([]string, 0, len(items))
	i := 1
	for _, item := range items {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d::int, $%d::int)", i, i+1))
		args = append(args, item.ProductID, item.Quantity)
		i += 2
	}

	query = fmt.Sprintf(query, strings.Join(valueStrings, ","))

	type result struct {
		ProductID   int     `db:"product_id"`
		ProductName string  `db:"product_name"`
		UnitPrice   float32 `db:"unit_price"`
		TotalPrice  float32 `db:"total_price"`
	}

	results := []result{}
	err := repo.sqlx.SelectContext(ctx, &results, query, args...)
	if err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return nil, err
	}

	orderItems := []*constant.OrderItems{}
	for _, result := range results {
		if item, ok := items[result.ProductID]; ok {
			item.ProductName = result.ProductName
			item.TotalPrice = result.TotalPrice
			item.UnitPrice = result.UnitPrice

			orderItems = append(orderItems, item)
		}
	}

	return orderItems, nil
}
