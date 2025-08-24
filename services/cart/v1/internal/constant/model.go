package constant

type Item struct {
	ProductID   int     `json:"product_id" bson:"product_id"`
	ProductName string  `json:"product_name" bson:"product_name"`
	Price       float64 `json:"price" bson:"price"`
	Quantity    int     `json:"quantity" bson:"quantity"`
	ImageURL    string  `json:"image_url" bson:"imageurl"`
}

type StoreItems struct {
	StoreID int     `json:"store_id" bson:"store_id"`
	Items   []*Item `json:"items" bson:"items"`
}

type Cart struct {
	UserID     string        `json:"user_id" bson:"user_id"`
	StoreItems []*StoreItems `json:"store_items" bson:"store_items"`
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
