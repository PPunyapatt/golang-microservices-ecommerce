package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"order/v1/internal/constant"
	"order/v1/internal/repository"
	"order/v1/proto/order"
	"package/rabbitmq"
	"package/rabbitmq/publisher"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type orderServer struct {
	tracer    trace.Tracer
	orderRepo repository.OrderRepository
	publisher publisher.EventPublisher
	order.UnimplementedOrderServiceServer
}

type orderService struct {
	tracer    trace.Tracer
	orderRepo repository.OrderRepository
	publisher publisher.EventPublisher
}

type OrderService interface {
	CreateProduct(context.Context, *constant.Product) error
	UpdateProduct(context.Context, *constant.Product) error
	UpdateStatus(context.Context, int, ...string) error
}

func NewOrderServer(orderRepo repository.OrderRepository, publisher publisher.EventPublisher, tracer trace.Tracer) (OrderService, order.OrderServiceServer) {
	return &orderService{
			tracer:    tracer,
			orderRepo: orderRepo,
			publisher: publisher,
		}, &orderServer{
			tracer:    tracer,
			orderRepo: orderRepo,
			publisher: publisher,
		}
}

func (s *orderServer) PlaceOrder(ctx context.Context, in *order.PlaceOrderRequest) (*order.PlaceOrderResponse, error) {
	orderCtx, orderSpan := s.tracer.Start(ctx, "create order")
	order := &constant.Order{
		UserID:          in.UserId,
		Status:          "pending",
		ShippingAddress: int(in.ShippingId),
	}

	orderItems := map[int]*constant.OrderItems{} //[]*constant.OrderItems{}
	for _, orderStore := range in.OrderItems {
		for _, item := range orderStore.Items {
			orderItems[int(item.ProductId)] = &constant.OrderItems{
				ProductID: int(item.ProductId),
				Quantity:  int(item.Quantity),
				StoreID:   int(orderStore.StoreId),
			}
		}
	}

	calCtx, calSpan := s.tracer.Start(orderCtx, "calculate price")
	items, err := s.orderRepo.CalculateTotalPrice(calCtx, orderItems)
	if err != nil {
		return nil, err
	}
	calSpan.End()

	total := float32(0)
	for _, item := range items {
		total += item.TotalPrice
	}

	order.TotalAmount = total

	oCtx, oSpan := s.tracer.Start(orderCtx, "create order and order items")
	tx, err := s.orderRepo.BeginTx()
	if err != nil {
		return nil, err
	}

	orderID, err := s.orderRepo.AddOrder(oCtx, tx, order)
	if err != nil {
		return nil, err
	}

	var invenOrder []*constant.InventoryOrder
	for _, item := range orderItems {
		item.OrderID = *orderID
		invenOrder = append(invenOrder, &constant.InventoryOrder{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	err = s.orderRepo.AddOrderItems(oCtx, tx, items)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	oSpan.End()
	orderSpan.End()

	payload := map[string]interface{}{
		"order_id":    order.OrderID,
		"total_price": order.TotalAmount,
		"items":       invenOrder,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	headers := amqp091.Table{}
	otel.GetTextMapPropagator().Inject(ctx, rabbitmq.AMQPHeaderCarrier(headers))

	if err = s.publisher.Publish(
		ctx,
		body,
		"order.exchange",
		"order.created",
		headers,
	); err != nil {
		return nil, err
	}

	log.Println("Publish success")

	return nil, nil
}

func (o *orderService) CreateProduct(ctx context.Context, product *constant.Product) error {
	// tracer := otel.Tracer("order-service")
	createCtx, createSpan := o.tracer.Start(ctx, "CreatedProduct")
	defer createSpan.End()
	if err := o.orderRepo.CreateProduct(createCtx, product); err != nil {
		return err
	}
	return nil
}

func (o *orderService) UpdateProduct(ctx context.Context, product *constant.Product) error {
	updateCtx, updateSpan := o.tracer.Start(ctx, "UpdateProduct")
	defer updateSpan.End()

	product.UpdatedAt = time.Now().UTC()
	if err := o.orderRepo.UpdateProduct(updateCtx, product); err != nil {
		return err
	}

	return nil
}

func (o *orderService) UpdateStatus(ctx context.Context, orderID int, args ...string) error {
	updateCtx, updateSpan := o.tracer.Start(ctx, fmt.Sprintf("update %s status", args[0]))
	defer updateSpan.End()
	if err := o.orderRepo.UpdateStatus(updateCtx, orderID, args[0], args[1]); err != nil {
		return err
	}
	return nil
}
