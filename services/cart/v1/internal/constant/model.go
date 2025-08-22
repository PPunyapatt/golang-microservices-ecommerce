package constant

import "time"

type Item struct {
	ProductID   int
	ProductName string
	Price       float64
	Quantity    int
	StoreID     int
	ImageURL    string
	UpdatedAt   *time.Time
	CreatedAt   *time.Time
}

type Cart struct {
	CartID    int
	UserID    string
	Items     []*Item
	UpdatedAt time.Time
	CreatedAt time.Time
}

type Pagination struct {
	Page       int32 `json:"page"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
	TotalCount int32 `json:"total_count"`
}

type PaymentData struct {
	UserID      string `json:"user_id"`
	OrderSource string `json:"order_source"`
}

func (Item) TableName() string {
	return "cart_items"
}

// func (Cart) TableName() string {
// 	return "cart"
// }
