package repository

import (
	"inventories/v1/internal/constant"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	AddInventory(*constant.Inventory) error
	UpdateInventory(*constant.Inventory) error
	RemoveInventory(int32) error
	GetInventory(int32) (*constant.Inventory, error)
	ListInventory(int32, string) ([]*constant.Inventory, error)

	AddCategory(*constant.Category) error
	UpdateCategory(*constant.Category) error
	RemoveCategory(int32) error
	GetCategories(string) (*constant.Category, error)
}

type inventoryRepository struct {
	gorm *gorm.DB
	sqlx *sqlx.DB
}

func NewInventoryRepository(gorm *gorm.DB, sqlx *sqlx.DB) InventoryRepository {
	return &inventoryRepository{
		gorm: gorm,
		sqlx: sqlx,
	}
}

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

func (repo *inventoryRepository) AddCategory(category *constant.Category) error {
	result := repo.gorm.Create(category)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *inventoryRepository) UpdateCategory(category *constant.Category) error {
	result := repo.gorm.Model(&constant.Category{}).Where("id = ?", category.ID).Updates(category)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *inventoryRepository) RemoveCategory(id int32) error {
	result := repo.gorm.Delete(&constant.Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *inventoryRepository) GetCategories(name string) (*constant.Category, error) {
	var category constant.Category
	query := repo.gorm.Model(&constant.Category{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	result := query.First(&category)
	if result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}
