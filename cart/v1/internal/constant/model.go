package constant

import "time"

type Item struct {
	CartID      int     `gorm:"column:cart_id"`
	ProductID   int     `gorm:"column:product_id"`
	ProductName string  `gorm:"column:name"`
	Price       float64 `gorm:"column:price"`
	Quantity    uint32  `gorm:"column:quantity"`
}

type Cart struct {
	CartID    int       `gorm:"column:id"`
	UserID    string    `gorm:"column:user_id"`
	Items     []*Item   `gorm:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (Item) TableName() string {
	return "cart_item"
}

func (Cart) TableName() string {
	return "cart"
}
