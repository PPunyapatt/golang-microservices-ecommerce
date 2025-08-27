package repository

import (
	"context"
	"fmt"
	"log"
	"order/v1/internal/constant"
	"order/v1/proto/order"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

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

func (repo *orderRepository) AddOrder(ctx context.Context, tx *gorm.DB, order *constant.Order, orderID *int) error {
	result := tx.WithContext(ctx).Create(order)
	if result.Error != nil {
		return result.Error
	}

	*orderID = order.OrderID

	return nil
}

func (repo *orderRepository) UpdateStatus(ctx context.Context, orderID int, data map[string]interface{}) error {
	result := repo.gorm.Model(&constant.Order{}).WithContext(ctx).Where("id = ?", orderID).Updates(data)
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

func (repo *orderRepository) GetItemsByOrderID(ctx context.Context, orderID int) ([]*constant.InventoryOrder, error) {
	var inventoryOrder []*constant.InventoryOrder
	result := repo.gorm.WithContext(ctx).
		Model(&constant.OrderItems{}).
		Where("order_id = ?", orderID).
		Select("product_id", "quantity").
		Find(&inventoryOrder)
	if result.Error != nil {
		return nil, result.Error
	}
	return inventoryOrder, nil
}

func (repo *orderRepository) CheckAndUpdateStatus(ctx context.Context, orderID int) error {
	result := repo.gorm.Model(&constant.Order{}).WithContext(ctx).Where("id = ? and (status = 'pending' OR status = 'reserved')", orderID).Update("status", "time_out")
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *orderRepository) ListOrder(ctx context.Context, req *constant.ListOrderRequest, pagination *constant.Pagination) ([]*order.Orders, error) {
	obj := repo.gorm.WithContext(ctx).Where("user_id = ?", req.UserID)
	if req.Status != nil {
		obj = obj.Where("sattus LIKE ?", "%"+*req.Status+"%")
	}

	// 2️⃣ Get total count (without limit/offset)
	var total int64
	if err := obj.Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.Total = int32(total)

	var orders []*constant.Order
	result := obj.Limit(int(pagination.Limit)).Offset(int(pagination.Offset)).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}

	if len(orders) == 0 {
		return []*order.Orders{}, nil
	}

	type OrderItemCount struct {
		OrderID int
		Count   int32
	}

	orderIDs := make([]int, len(orders))
	for i, o := range orders {
		orderIDs[i] = o.OrderID
	}

	var counts []OrderItemCount
	if err := repo.gorm.WithContext(ctx).
		Model(&constant.OrderItems{}).
		Select("order_id, COUNT(*) as count").
		Where("order_id IN ?", orderIDs).
		Group("order_id").
		Scan(&counts).Error; err != nil {
		return nil, err
	}

	countMap := make(map[int]int32, len(counts))
	for orderID, count := range countMap {
		countMap[orderID] = count
	}

	ordersRPC := make([]*order.Orders, len(orders))

	for _, order_ := range orders {
		count := countMap[order_.OrderID]

		ordersRPC = append(ordersRPC, &order.Orders{
			OrderId:       int32(order_.OrderID),
			TotalAmount:   float64(order_.TotalAmount),
			Status:        order_.Status,
			PaymentStatus: order_.PaymentStatus,
			TotalItems:    count,
		})
	}

	return ordersRPC, nil
}

func (repo *orderRepository) CheckOrderStatus(ctx context.Context, orderID int, status ...string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM orders
			WHERE id = ? AND status IN (?)
		)
	`

	args := []interface{}{orderID, status}
	query, args, err := sqlx.In(query, args...)
	if err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return false, err
	}

	query = repo.sqlx.Rebind(query)

	var exists bool
	if err := repo.sqlx.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return false, nil
	}

	return exists, nil
}
