package service

import (
	"context"
	"order/v1/internal/repository"
	"order/v1/proto/order"
)

type orderServer struct {
	orderRepo repository.OrderRepository
	order.UnimplementedOrderServiceServer
}

func NewOrderServer(orderRepo repository.OrderRepository) order.OrderServiceServer {
	return &orderServer{orderRepo: orderRepo}
}

func (s *orderServer) PlaceOrder(ctx context.Context, in *order.PlaceOrderRequest) (*order.PlaceOrderResponse, error) {

	return nil, nil
}
