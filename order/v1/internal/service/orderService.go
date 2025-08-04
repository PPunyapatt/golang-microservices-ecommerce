package service

import (
	"context"
	"order/v1/internal/repository"
	"order/v1/pkg/rabbitmq/publisher"
	"order/v1/proto/order"
)

type orderServer struct {
	orderRepo repository.OrderRepository
	publisher publisher.EventPublisher
	order.UnimplementedOrderServiceServer
}

func NewOrderServer(orderRepo repository.OrderRepository, publisher publisher.EventPublisher) order.OrderServiceServer {
	return &orderServer{
		orderRepo: orderRepo,
		publisher: publisher,
	}
}

func (s *orderServer) PlaceOrder(ctx context.Context, in *order.PlaceOrderRequest) (*order.PlaceOrderResponse, error) {

	return nil, nil
}
