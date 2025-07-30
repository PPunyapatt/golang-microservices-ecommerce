package repository

import (
	"inventories/v1/internal/constant"
	"inventories/v1/proto/Inventory"

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
