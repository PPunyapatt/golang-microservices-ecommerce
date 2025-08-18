package constant

import "time"

type Order struct {
	OrderID         int       `json:"order_id" db:"order_id" gorm:"column:id;primaryKey"`
	UserID          string    `json:"user_id" db:"user_id"`
	Status          string    `json:"status" db:"status"`
	TotalAmount     float32   `json:"total_amount" db:"total_amount"`
	PaymentID       int       `json:"payment_id" db:"payment_id"`
	PaymentStatus   string    `json:"payment_status" db:"payment_status"`
	ShippingAddress int       `json:"shipping_address_id" db:"shipping_address_id" gorm:"column:shipping_address_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type OrderItems struct {
	OrderItemID int     `json:"order_item_id" db:"order_item_id" gorm:"column:id;primaryKey"`
	OrderID     *int    `json:"order_id" db:"order_id"`
	ProductID   int     `json:"product_id" db:"product_id"`
	Quantity    int     `json:"quantity" db:"quantity"`
	TotalPrice  float32 `json:"total_price" db:"total_price"`
	StoreID     int     `json:"store_id" db:"store_id"`
	ProductName string  `json:"product_name" db:"product_name"`
	UnitPrice   float32 `json:"unit_price" db:"unit_price"`
}

type UpdateStatus struct {
	OrderID int    `json:"order_id"`
	Status  string `json:"status"`
}

type Product struct {
	StoreID     *int      `json:"store_id"`
	ProductID   *int      `json:"product_id"`
	ProductName *string   `json:"product_name"`
	Price       *float32  `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type InventoryOrder struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

func (Product) TableName() string {
	return "order_products"
}
