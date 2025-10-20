package handler

import (
	"gateway/v1/internal/constant"
	"gateway/v1/proto/Inventory"
	"gateway/v1/proto/auth"
	"gateway/v1/proto/cart"
	"gateway/v1/proto/order"
	"gateway/v1/proto/payment"
	"package/metrics"
)

type ApiHandler struct {
	CartSvc           cart.CartServiceClient
	AuthSvc           auth.AuthServiceClient
	InventorySvc      Inventory.InventoryServiceClient
	PaymentSvc        payment.PaymentServiceClient
	OrderSvc          order.OrderServiceClient
	prometheusMetrics *metrics.Metrics
}

func ServiceNew(svc *constant.Clients, prometheusMetrics *metrics.Metrics) *ApiHandler {
	return &ApiHandler{
		AuthSvc:           svc.AuthClient,
		CartSvc:           svc.CartClient,
		InventorySvc:      svc.InventoryClient,
		PaymentSvc:        svc.PaymentClient,
		OrderSvc:          svc.OrderClient,
		prometheusMetrics: prometheusMetrics,
	}
}
