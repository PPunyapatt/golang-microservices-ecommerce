package constant

type PlaceOrderRequest struct {
	UserID     string
	ShippingID int32              `json:"shipping_id" validate:"required"`
	OrderItems []*PlaceOrderStore `json:"order_items" validate:"required,dive"`
}

type PlaceOrderStore struct {
	StoreID int32              `json:"store_id" validate:"required"`
	Items   []*PlaceOrderItems `json:"items" validate:"required,dive"`
}

type PlaceOrderItems struct {
	OrderID   int32 `json:"order_id"`
	ProductID int32 `json:"product_id"  validate:"required"`
	Quantity  int32 `json:"quantity"  validate:"required"`
}
