package constant

import "time"

type Pagination struct {
	Limit  int32 `json:"page" db:"page"`
	Offset int32 `json:"offset" db:"offset"`
	Total  int32 `json:"total" db:"total"`
}

type Inventory struct {
	ID             int32     `json:"id" db:"id"` // ID ยังใช้ int32 ปกติ
	StoreID        *int32    `json:"store_id" db:"store_id"`
	AddBy          *string   `json:"add_by" db:"add_by"`
	Name           *string   `json:"name" db:"name"`
	Description    *string   `json:"description" db:"description"`
	Price          *float64  `json:"price" db:"price"`
	AvailableStock *int32    `json:"available_stock" db:"available_stock"`
	ReservedStock  *int32    `json:"reserved_stock" db:"reserved_stock"`
	CategoryID     *int32    `json:"category_id" db:"category_id"`
	ImageURL       *string   `json:"image_url" db:"image_url"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type Category struct {
	ID      int32  `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	StoreID int32  `json:"store_id" db:"store_id"`
}

type ListInventoryReq struct {
	StoreID    *int32  `json:"store_id" db:"store_id"`
	Query      *string `json:"search" db:"search"`
	CategoryID *int32  `json:"category_id" db:"category_id"`
}

func (Inventory) TableName() string {
	return "products"
}
