package constant

type Items struct {
	ProductID   int32   `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int32   `json:"quantity"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
}

type StoreItems struct {
	StoreID int32    `json:"store_id"`
	Items   []*Items `json:"items" validate:"required,dive"`
}

type Products struct {
	Products []*StoreItems `json:"products"`
}

type Cart struct {
}

type RemoveItemReq struct {
	CartItemID int32 `params:"cart_item_id"`
	CartID     int32 `params:"cart_id"`
}
