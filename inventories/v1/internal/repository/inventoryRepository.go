package repository

import (
	"inventories/v1/internal/constant"
	"inventories/v1/proto/Inventory"
	"log"

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
	ListCategories(int32, string, *constant.Pagination) ([]*Inventory.Category, error)
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

func (repo *inventoryRepository) ListCategories(store_id int32, name string, pagination *constant.Pagination) ([]*Inventory.Category, error) {
	var categories []*constant.Category
	query := repo.gorm.
		Model(&constant.Category{}).
		Where("store_id = ? AND name LIKE ?", int(store_id), "%"+name+"%")

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.Total = int32(total)
	log.Println("Total categories found:", pagination.Total)

	if pagination.Offset == 0 {
		query = query.
			Limit(int(pagination.Limit)).
			Offset(int(pagination.Offset))
	}

	if err := query.Find(&categories).Error; err != nil {
		return nil, err
	}

	result := []*Inventory.Category{}
	for _, cat := range categories {
		result = append(result, &Inventory.Category{
			CategoryID: cat.ID,
			Name:       cat.Name,
			StoreID:    cat.StoreID,
		})
	}
	return result, nil
}
