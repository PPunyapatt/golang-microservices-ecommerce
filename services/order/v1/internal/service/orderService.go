package service

import (
	"context"
	"encoding/json"
	"order/v1/internal/repository"
	"order/v1/proto/order"
	"package/rabbitmq/publisher"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *orderServer) UpdateStatus(ctx context.Context, in *order.UpdateStatusRequest) (*order.Empty, error) {
	payload := map[string]interface{}{
		"order_id": in.OrderID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to marshal payload")
	}

	if err = s.publisher.Publish(ctx, body); err != nil {
		return nil, err
	}
	return nil, nil
}
