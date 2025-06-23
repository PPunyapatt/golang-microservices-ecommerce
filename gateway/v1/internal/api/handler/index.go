package handler

import (
	"gateway/v1/proto/auth"
	"gateway/v1/proto/cart"
)

type ApiHandler struct {
	CartSvc cart.CartServiceClient
	AuthSvc auth.AuthServiceClient
}

func ServiceNew(auth auth.AuthServiceClient, cart cart.CartServiceClient) *ApiHandler {
	return &ApiHandler{
		AuthSvc: auth,
		CartSvc: cart,
	}
}
