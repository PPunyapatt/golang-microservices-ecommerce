package repository

import (
	"inventories/v1/internal/constant"
)

func (repo *inventoryRepository) AddInventory(inventory *constant.Inventory) error {
	result := repo.gorm.Create(inventory)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *inventoryRepository) UpdateInventory(in *constant.Inventory) error {

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
	if in.Stock != nil {
		updateData["stock"] = in.Stock
	}
	if in.CategoryID != nil {
		updateData["category_id"] = in.CategoryID
	}
	if in.ImageURL != nil {
		updateData["image_url"] = in.ImageURL
	}

	result := repo.gorm.Model(&constant.Inventory{}).Where("id = ?", in.ID).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *inventoryRepository) RemoveInventory(id int32) error {
	result := repo.gorm.Delete(&constant.Inventory{}, id)
	if result.Error != nil {
		return result.Error
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

func (repo *inventoryRepository) ListInventory(storeID int32, name string) ([]*constant.Inventory, error) {
	var inventories []*constant.Inventory
	query := repo.gorm.Model(&constant.Inventory{})

	if storeID > 0 {
		query = query.Where("store_id = ?", storeID)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	result := query.Find(&inventories)
	if result.Error != nil {
		return nil, result.Error
	}
	return inventories, nil
}
