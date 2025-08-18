package repository

import (
	"context"
	"inventories/v1/internal/constant"
	"inventories/v1/proto/Inventory"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	// ---- Inventory ----
	AddInventory(context.Context, *constant.Inventory) (*int32, error)
	UpdateInventory(context.Context, *constant.Inventory) error
	RemoveInventory(userID string, storeID, Inventory int32) error
	GetInventory(int32) (*constant.Inventory, error)
	ListInventory(*constant.ListInventoryReq, *constant.Pagination) ([]*Inventory.Inventory, error)
	ReserveStock(context.Context, []*constant.Item) error
	CutStock(context.Context, []*constant.Item) error
	ReleaseStock(context.Context, []*constant.Item) error

	// ---- Category ----
	AddCategory(*constant.Category) error
	UpdateCategory(*constant.Category) error
	RemoveCategory(int32) error
	ListCategories(*constant.ListInventoryReq, *constant.Pagination) ([]*Inventory.Category, error)
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
