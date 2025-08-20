package service

import (
	"context"
	"encoding/json"
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
	UpdateStatus(context.Context, int, map[string]interface{}) error

	PushEventCutorReleaseStock(context.Context, int, string) error
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

	orderID := 0                                 // Temporary
	orderItems := map[int]*constant.OrderItems{} //[]*constant.OrderItems{}
	for _, orderStore := range in.OrderItems {
		for _, item := range orderStore.Items {
			orderItems[int(item.ProductId)] = &constant.OrderItems{
				OrderID:   &orderID,
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

	err = s.orderRepo.AddOrder(oCtx, tx, order, &orderID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var invenOrder []*constant.InventoryOrder
	for _, item := range orderItems {
		invenOrder = append(invenOrder, &constant.InventoryOrder{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	err = s.orderRepo.AddOrderItems(oCtx, tx, items)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	oSpan.End()
	orderSpan.End()

	payload := map[string]interface{}{
		"order_id":    orderID,
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
		1,
	); err != nil {
		return nil, err
	}

	if err = s.publisher.Publish(
		ctx,
		body,
		"order.exchange",
		"order.dlq",
		headers,
		1,
	); err != nil {
		return nil, err
	}

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

func (o *orderService) UpdateStatus(ctx context.Context, orderID int, args map[string]interface{}) error {
	updateCtx, updateSpan := o.tracer.Start(ctx, "update status")
	defer updateSpan.End()
	log.Println("args: ", args)
	if err := o.orderRepo.UpdateStatus(updateCtx, orderID, args); err != nil {
		return err
	}
	return nil
}

func (o *orderService) PushEventCutorReleaseStock(ctx context.Context, orderID int, key string) error {
	inventoryOrder, err := o.orderRepo.GetItemsByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	body, err := json.Marshal(inventoryOrder)
	if err != nil {
		return err
	}

	var routingKey string
	switch key {
	case "payment.seccussed":
		routingKey = "order.payment.successed"
	case "payment.failed":
		routingKey = "order.payment.failed"
	}

	headers := amqp091.Table{}
	otel.GetTextMapPropagator().Inject(ctx, rabbitmq.AMQPHeaderCarrier(headers))

	if err = o.publisher.Publish(
		ctx,
		body,
		"order.exchange",
		routingKey,
		headers,
	); err != nil {
		return err
	}
	return nil
}
