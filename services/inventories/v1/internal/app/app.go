package app

import (
	"context"
	"encoding/json"
	"inventories/v1/internal/constant"
	"inventories/v1/internal/services"
	"log"
	"log/slog"
	"package/rabbitmq"

	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

type AppServer interface {
	Worker(ctx context.Context, messages <-chan amqp091.Delivery)
}

type appServer struct {
	inventoryService services.InventoryServie
}

func NewWorker(inventoryService services.InventoryServie) AppServer {
	return &appServer{
		inventoryService: inventoryService,
	}
}

func (c *appServer) Worker(ctx context.Context, messages <-chan amqp091.Delivery) {
	for delivery := range messages {
		log.Println("delivery.Type: ", delivery.RoutingKey)
		switch delivery.RoutingKey {
		case "order.created":
			c.ReserveStock(ctx, delivery)
		case "order.payment.successed":
			c.CutOrReleaseStock(ctx, delivery, delivery.RoutingKey)
		case "order.payment.failed":
			c.CutOrReleaseStock(ctx, delivery, delivery.RoutingKey)
		case "order.timeout":
			c.CutOrReleaseStock(ctx, delivery, delivery.RoutingKey)
		}
	}
}

func (c *appServer) ReserveStock(ctx context.Context, delivery amqp091.Delivery) {
	ctx_ := otel.GetTextMapPropagator().Extract(
		ctx,
		rabbitmq.AMQPHeaderCarrier(delivery.Headers),
	)

	var payload constant.Order
	err := json.Unmarshal(delivery.Body, &payload)
	if err != nil {
		slog.Error("failed to Unmarshal", err)
	}
	if err = c.inventoryService.ReserveStock(ctx_, &payload); err != nil {
		slog.Error("failed to reserve stock", err)
		c.rejectDelivery(delivery)
		return
	}

	c.ackDelivery(delivery)
}

func (c *appServer) CutOrReleaseStock(ctx context.Context, delivery amqp091.Delivery, routingKey string) {
	var payload []*constant.Item
	err := json.Unmarshal(delivery.Body, &payload)
	if err != nil {
		slog.Error("failed to Unmarshal", err)
	}

	ctx_ := otel.GetTextMapPropagator().Extract(
		ctx,
		rabbitmq.AMQPHeaderCarrier(delivery.Headers),
	)

	var handler func(context.Context, []*constant.Item) error

	switch routingKey {
	case "order.payment.successed":
		handler = c.inventoryService.CutStock
	case "order.payment.failed":
		handler = c.inventoryService.ReleaseStock
	case "order.timeout":
		handler = c.inventoryService.ReleaseStock
	}

	if err := handler(ctx_, payload); err != nil {
		slog.Error("failed to cut stock", err)
		c.rejectDelivery(delivery)
		return
	}

	c.ackDelivery(delivery)
}

// -------------------------- Handler Error --------------------------
func (c *appServer) handleUnknownDelivery(delivery amqp091.Delivery) {
	slog.Warn("unknown delivery routing key", "key", delivery.RoutingKey)
	c.rejectDelivery(delivery)
}

func (c *appServer) rejectDelivery(delivery amqp091.Delivery) {
	if err := delivery.Reject(false); err != nil {
		slog.Error("failed to delivery.Reject", err)
	}
}

func (c *appServer) ackDelivery(delivery amqp091.Delivery) {
	err := delivery.Ack(false)
	if err != nil {
		slog.Error("failed to acknowledge delivery", err)
	} else {
		slog.Info("ack success", "delivery_tag", delivery.DeliveryTag)
	}
}
