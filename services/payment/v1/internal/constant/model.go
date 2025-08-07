package constant

type PaymentRequest struct {
	OrderID     int     `json:"order_id"`
	AmountPrice float32 `json:"amount_price"`
	UserID      string  `json:"user_id"`
}
