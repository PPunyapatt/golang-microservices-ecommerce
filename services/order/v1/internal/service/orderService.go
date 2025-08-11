package service

import (
	"context"
	"encoding/json"
	"order/v1/internal/constant"
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

	tx, err := s.orderRepo.BeginTx()
	if err != nil {
		return nil, err
	}

	// amount := 0
	for _, i := range in.OrderItems {
		amount += (i.Quantity * i)
	}

	// order := &constant.Order{
	// 	UserID: in.UserId,
	// 	Status: "pending",
	// 	TotalAmount: in.OrderItems.,
	// }

	// orderID, err := s.orderRepo.AddOrder(tx, ctx)
	// if err != nil {
	// 	return nil, err
	// }

	var orderID *int
	var total_amount float32
	orderItems := []*constant.OrderItems{}
	for _, orderItem := range in.OrderItems {
		for _, item := range orderItem.Items {
			itemData := &constant.OrderItems{
				OrderID:    *orderID,
				StoreID:    int(orderItem.StoreId),
				ProductID:  int(item.ProductId),
				UnitPrice:  float32(item.UnitPrice),
				Quantity:   int(item.Quantity),
				TotalPrice: float32(float64(item.Quantity) * item.UnitPrice),
			}
		}

	}

	return nil, nil
}

func (s *orderServer) OrderID(orderID int) int32 {
	return int32(orderID)
}

func (s *orderServer) UpdateStatus(ctx context.Context, in *order.UpdateStatusRequest) (*order.Empty, error) {
	payload := map[string]interface{}{
		"order_id": in.OrderID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to marshal payload")
	}

	if err = s.publisher.Publish(ctx, body, "order.update"); err != nil {
		return nil, err
	}
	return nil, nil
}
