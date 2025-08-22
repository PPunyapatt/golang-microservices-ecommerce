package constant

import (
	"time"
)

type PaymentRequest struct {
	OrderID     int32   `json:"order_id"`
	TotalPrice  float32 `json:"total_price"`
	UserID      string  `json:"user_id"`
	OrderSource string  `json:"order_source"`
	// Currency      string  `json:"currency"`
	// PaymentMethod string  `json:"payment_method"`
}

type Payment struct {
	PaymentID     string    `json:"payment_id"`
	OrderID       int       `json:"order_id"`
	Amount        float32   `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	FailureReason *string   `json:"failure_reason"`
}
