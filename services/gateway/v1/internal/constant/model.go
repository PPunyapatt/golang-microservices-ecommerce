package constant

import (
	"gateway/v1/proto/Inventory"
	"gateway/v1/proto/auth"
	"gateway/v1/proto/cart"
	"gateway/v1/proto/payment"
)

type Clients struct {
	CartClient      cart.CartServiceClient
	AuthClient      auth.AuthServiceClient
	InventoryClient Inventory.InventoryServiceClient
	PaymentClient   payment.PaymentServiceClient
}
