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

	"package/metrics"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type orderServer struct {
	tracer    trace.Tracer
	orderRepo repository.OrderRepository
	publisher publisher.EventPublisher
	pm        *metrics.Metrics
	order.UnimplementedOrderServiceServer
}

type orderService struct {
	tracer    trace.Tracer
	orderRepo repository.OrderRepository
	publisher publisher.EventPublisher
	pm        *metrics.Metrics
}

type OrderService interface {
	ProductUpdate(context.Context, *constant.OrderProduct) error
	UpdateStatus(context.Context, int, map[string]interface{}) error
	CheckAndUpdateStatus(context.Context, int) error

	PushEventCutorReleaseStock(context.Context, int, string) error
}

func NewOrderServer(orderRepo repository.OrderRepository, publisher publisher.EventPublisher, tracer trace.Tracer, pm *metrics.Metrics) (OrderService, order.OrderServiceServer) {
	return &orderService{
			tracer:    tracer,
			orderRepo: orderRepo,
			publisher: publisher,
			pm:        pm,
		}, &orderServer{
			tracer:    tracer,
			orderRepo: orderRepo,
			publisher: publisher,
			pm:        pm,
		}
}

func (s *orderServer) PlaceOrder(ctx context.Context, in *order.PlaceOrderRequest) (*order.Empty, error) {
	s.pm.Grpc.OrderPlaceRequests.Inc()
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
		log.Printf("%+v", errors.WithStack(err))
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
		log.Printf("%+v", errors.WithStack(err))
		return nil, err
	}

	err = s.orderRepo.AddOrder(oCtx, tx, order, &orderID)
	if err != nil {
		tx.Rollback()
		log.Printf("%+v", errors.WithStack(err))
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
		log.Printf("%+v", errors.WithStack(err))
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return nil, err
	}
	oSpan.End()
	orderSpan.End()

	payload := map[string]interface{}{
		"user_id":      in.UserId,
		"order_id":     orderID,
		"total_price":  order.TotalAmount,
		"items":        invenOrder,
		"order_source": in.OrderSource,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("%+v", errors.WithStack(err))
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
		log.Printf("%+v", errors.WithStack(err))
		return nil, err
	}

	return nil, nil
}

func (s *orderServer) ListOrder(ctx context.Context, in *order.ListOrderRequest) (*order.ListOrderResponse, error) {
	listCtx, listSpan := s.tracer.Start(ctx, "CreatedProduct")
	req := &constant.ListOrderRequest{
		UserID: in.UserId,
		Status: in.Status,
	}
	pagination := &constant.Pagination{
		Limit:  in.Pagination.Limit,
		Offset: in.Pagination.Offset,
	}
	orders, err := s.orderRepo.ListOrder(listCtx, req, pagination)
	if err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return nil, err
	}
	listSpan.End()

	result := &order.ListOrderResponse{
		Orders: orders,
		Pagination: &order.Pagination{
			Limit: pagination.Limit,
			Total: &pagination.Total,
		},
	}
	return result, nil
}

func (o *orderService) ProductUpdate(ctx context.Context, orderProduct *constant.OrderProduct) error {
	product := &constant.Product{
		StoreID:     &orderProduct.StoreID,
		ProductID:   &orderProduct.ProductID,
		ProductName: &orderProduct.ProductName,
		Price:       &orderProduct.Price,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	switch orderProduct.Operation {
	case "c":
		if err := o.orderRepo.CreateProduct(ctx, product); err != nil {
			log.Printf("%+v", errors.WithStack(err))
			return err
		}
	case "u":
		if err := o.orderRepo.UpdateProduct(ctx, product); err != nil {
			log.Printf("%+v", errors.WithStack(err))
			return err
		}
	case "d":
		if err := o.orderRepo.DeleteProduct(ctx, *product.ProductID); err != nil {
			log.Printf("%+v", errors.WithStack(err))
			return err
		}
	}
	return nil
}

func (o *orderService) UpdateStatus(ctx context.Context, orderID int, args map[string]interface{}) error {
	updateCtx, updateSpan := o.tracer.Start(ctx, "update status")
	defer updateSpan.End()
	if err := o.orderRepo.UpdateStatus(updateCtx, orderID, args); err != nil {
		return err
	}
	return nil
}

func (o *orderService) PushEventCutorReleaseStock(ctx context.Context, orderID int, key string) error {
	inventoryOrder, err := o.orderRepo.GetItemsByOrderID(ctx, orderID)
	if err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return err
	}

	body, err := json.Marshal(inventoryOrder)
	if err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return err
	}

	var routingKey string
	switch key {
	case "payment.successed":
		routingKey = "order.payment.successed"
	case "payment.failed":
		routingKey = "order.payment.failed"
	case "order.timeout":
		routingKey = key
		exist, err := o.orderRepo.CheckOrderStatus(ctx, orderID, "reserved", "pending")
		if err != nil {
			log.Printf("%+v", errors.WithStack(err))
			return err
		}

		if !exist {
			log.Println("Order status isn't pending or reserved")
			return nil
		}
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
		log.Printf("%+v", errors.WithStack(err))
		return err
	}
	return nil
}

func (o *orderService) CheckAndUpdateStatus(ctx context.Context, orderID int) error {
	err := o.orderRepo.CheckAndUpdateStatus(ctx, orderID)
	if err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return err
	}
	return nil
}
