package handler

import (
	"gateway/v1/internal/constant"
	"gateway/v1/proto/Inventory"
	"gateway/v1/proto/auth"
	"gateway/v1/proto/cart"
)

type ApiHandler struct {
	CartSvc      cart.CartServiceClient
	AuthSvc      auth.AuthServiceClient
	InventorySvc Inventory.InventoryServiceClient
}

func ServiceNew(
	svc *constant.Clients,
) *ApiHandler {
	return &ApiHandler{
		AuthSvc:      svc.AuthClient,
		CartSvc:      svc.CartClient,
		InventorySvc: svc.InventoryClient,
	}
}
