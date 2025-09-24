package repository

import (
	"context"
	"errors"
	"fmt"
	"inventories/v1/internal/constant"
	"inventories/v1/proto/Inventory"
	"log"
	"strings"
)

func (repo *inventoryRepository) AddInventory(ctx context.Context, inventory *constant.Inventory) (*int32, error) {
	result := repo.gorm.WithContext(ctx).Omit("updated_at").Create(inventory)
	if result.Error != nil {
		return nil, result.Error
	}
	return &inventory.ID, nil
}

func (repo *inventoryRepository) UpdateInventory(ctx context.Context, in *constant.Inventory) error {

	updateData := map[string]interface{}{}

	if in.StoreID != nil {
		updateData["store_id"] = in.StoreID
	}
	if in.AddBy != nil {
		updateData["add_by"] = in.AddBy
	}
	if in.Name != nil {
		updateData["name"] = in.Name
	}
	if in.Description != nil {
		updateData["description"] = in.Description
	}
	if in.Price != nil {
		updateData["price"] = in.Price
	}
	if in.AvailableStock != nil {
		updateData["available_stock"] = in.AvailableStock
	}
	if in.CategoryID != nil {
		updateData["category_id"] = in.CategoryID
	}
	if in.ImageURL != nil {
		updateData["image_url"] = in.ImageURL
	}

	log.Println("Update Inventory Data:", updateData)

	result := repo.gorm.Model(&constant.Inventory{}).WithContext(ctx).Where("id = ?", in.ID).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *inventoryRepository) RemoveInventory(ctx context.Context, userID string, storeID, inventoryID int32) error {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM stores s
			LEFT JOIN products p
				ON s.id = p.store_id
			WHERE s.owner = $1 AND s.id = $2
		)
	`

	args := []interface{}{userID, storeID}
	var exists bool
	err := repo.sqlx.DB.QueryRow(query, args...).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user is not authorized")
	}

	result := repo.gorm.Delete(&constant.Inventory{}, inventoryID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("Item does not exist!!!")
	}

	return nil
}

func (repo *inventoryRepository) GetInventory(id int32) (*constant.Inventory, error) {
	var inventory constant.Inventory
	result := repo.gorm.First(&inventory, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &inventory, nil
}

func (repo *inventoryRepository) ListInventory(ctx context.Context, req *constant.ListInventoryReq, pagination *constant.Pagination) ([]*Inventory.Inventory, error) {
	var inventories []*constant.Inventory
	query := repo.gorm.Model(&constant.Inventory{})

	if req.StoreID != nil {
		query = query.Where("store_id = ?", *req.StoreID)
	}
	if req.Query != nil {
		query = query.Where("name LIKE ?", "%"+*req.Query+"%")
	}
	if req.CategoryID != nil {
		query = query.Where("category_id = ?", *req.CategoryID)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.Total = int32(total)

	err := query.Find(&inventories)
	if err.Error != nil {
		return nil, err.Error
	}

	result := []*Inventory.Inventory{}
	for _, inv := range inventories {
		result = append(result, &Inventory.Inventory{
			ID:          &inv.ID,
			Name:        inv.Name,
			Description: inv.Description,
			Price:       inv.Price,
			Stock:       inv.AvailableStock,
			CategoryID:  inv.CategoryID,
			ImageURL:    inv.ImageURL,
			StoreID:     inv.StoreID,
			AddBy:       inv.AddBy,
		})
	}
	return result, nil
}

func (repo *inventoryRepository) ReserveStock(ctx context.Context, items []*constant.Item) error {
	query := `
		WITH insufficient AS (
			SELECT
				p.id as product_id,
				p.available_stock,
				i.quantity
			FROM products p
			INNER JOIN (
				VALUES
					%s
			) i(product_id, quantity)
				ON i.product_id = p.id 
			WHERE p.available_stock < i.quantity
		)
		UPDATE products p
		SET 
			available_stock = p.available_stock - i.quantity, 
			reserved_stock = i.quantity
		FROM (
			VALUES
				%s
		) i(product_id, quantity)
		WHERE
			p.id = i.product_id AND 
			NOT EXISTS (SELECT 1 FROM insufficient)
	`

	valueStrings := make([]string, 0, len(items))
	for _, order := range items {
		orderValue := fmt.Sprintf("(%d::int, %d::int)", order.ProductID, order.Quantity)
		valueStrings = append(valueStrings, orderValue)
	}

	values := strings.Join(valueStrings, ", ")
	query = fmt.Sprintf(query, values, values)

	result, err := repo.sqlx.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return errors.New("The product is out of stock.")
	}

	return nil
}

func (repo *inventoryRepository) CutStock(ctx context.Context, items []*constant.Item) error {
	query := `
		UPDATE products p
		SET 
			reserved_stock = p.reserved_stock - i.quantity,
			updated_at = NOW()
		FROM (
			VALUES
				%s
		) i(product_id, quantity)
		WHERE p.id = i.product_id
	`

	valueStrings := make([]string, 0, len(items))
	args := []interface{}{}
	for i, item := range items {
		n := i*2 + 1
		valueStrings = append(valueStrings, fmt.Sprintf("($%d::int, $%d::int)", n, n+1))
		args = append(args, item.ProductID, item.Quantity)
	}

	query = fmt.Sprintf(query, strings.Join(valueStrings, ", "))
	_, err := repo.sqlx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repo *inventoryRepository) ReleaseStock(ctx context.Context, items []*constant.Item) error {
	query := `
		UPDATE products p
		SET 
			available_stock = p.available_stock + i.quantity,
            reserved_stock = p.reserved_stock - i.quantity,
			updated_at = NOW()
		FROM (
			VALUES
				%s
		) i(product_id, quantity)
		WHERE p.id = i.product_id
	`

	valueStrings := make([]string, 0, len(items))
	args := []interface{}{}
	for i, item := range items {
		n := i*2 + 1
		valueStrings = append(valueStrings, fmt.Sprintf("($%d::int, $%d::int)", n, n+1))
		args = append(args, item.ProductID, item.Quantity)
	}

	query = fmt.Sprintf(query, strings.Join(valueStrings, ", "))
	_, err := repo.sqlx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
