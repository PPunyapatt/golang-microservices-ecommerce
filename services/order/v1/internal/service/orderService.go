package service

import (
	"context"
	"order/v1/internal/constant"
	"order/v1/internal/repository"
	"order/v1/proto/order"
	"package/rabbitmq/publisher"
	"time"

	"go.opentelemetry.io/otel"
)

type orderServer struct {
	orderRepo repository.OrderRepository
	publisher publisher.EventPublisher
	order.UnimplementedOrderServiceServer
}

type orderService struct {
	orderRepo repository.OrderRepository
	publisher publisher.EventPublisher
}

type OrderService interface {
	CreateProduct(context.Context, *constant.Product) error
	UpdateProduct(context.Context, *constant.Product) error
	UpdateStatus(context.Context) error
}

func NewOrderServer(orderRepo repository.OrderRepository, publisher publisher.EventPublisher) (OrderService, order.OrderServiceServer) {
	return &orderService{
			orderRepo: orderRepo,
			publisher: publisher,
		}, &orderServer{
			orderRepo: orderRepo,
			publisher: publisher,
		}
}

func (s *orderServer) PlaceOrder(ctx context.Context, in *order.PlaceOrderRequest) (*order.PlaceOrderResponse, error) {

	// tx, err := s.orderRepo.BeginTx()
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

func (s *orderServer) OrderID(orderID int) int32 {
	return int32(orderID)
}

func (o *orderService) CreateProduct(ctx context.Context, product *constant.Product) error {
	tracer := otel.Tracer("order-service")
	createCtx, createSpan := tracer.Start(ctx, "CreatedProduct")
	defer createSpan.End()
	if err := o.orderRepo.CreateProduct(createCtx, product); err != nil {
		return err
	}
	return nil
}

func (o *orderService) UpdateProduct(ctx context.Context, product *constant.Product) error {
	tracer := otel.Tracer("order-service")
	updateCtx, updateSpan := tracer.Start(ctx, "UpdateProduct")
	defer updateSpan.End()
	product.UpdatedAt = time.Now().UTC()
	if err := o.orderRepo.UpdateProduct(updateCtx, product); err != nil {
		return err
	}

	return nil
}

func (o *orderService) UpdateStatus(ctx context.Context) error {
	return nil
}
