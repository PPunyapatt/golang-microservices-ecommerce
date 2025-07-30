package repository

import (
	"inventories/v1/internal/constant"
	"inventories/v1/proto/Inventory"
	"log"
)

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
