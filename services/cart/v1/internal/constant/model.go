package constant

import "time"

type Item struct {
	CartID      int        `gorm:"column:cart_id"`
	ProductID   int        `gorm:"column:product_id" db:"product_id"`
	ProductName string     `gorm:"column:product_name" db:"product_name"`
	Price       float64    `gorm:"column:price" db:"price"`
	Quantity    int        `gorm:"column:quantity" db:"quantity"`
	StoreID     int        `gorm:"column:store_id" db:"store_id"`
	ImageURL    string     `gorm:"column:image_url" db:"image_url"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
	CreatedAt   *time.Time `gorm:"column:created_at"`
}

type Cart struct {
	CartID    int       `gorm:"column:id"`
	UserID    string    `gorm:"column:user_id"`
	Items     []*Item   `gorm:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

type Pagination struct {
	Page       int32 `json:"page"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
	TotalCount int32 `json:"total_count"`
}

func (Item) TableName() string {
	return "cart_items"
}

// func (Cart) TableName() string {
// 	return "cart"
// }
