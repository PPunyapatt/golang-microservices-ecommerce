package app

import (
	"context"
	"encoding/json"
	"log/slog"
	"order/v1/internal/constant"
	"order/v1/internal/service"
	"package/rabbitmq"

	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

type AppServer interface {
	Worker(ctx context.Context, messages <-chan amqp091.Delivery)
}

type appServer struct {
	orderService service.OrderService
}

func NewWorker(orderService service.OrderService) AppServer {
	return &appServer{
		orderService: orderService,
	}
}

func (c *appServer) Worker(ctx context.Context, messages <-chan amqp091.Delivery) {
	for delivery := range messages {
		slog.Info("processDeliveries", "delivery_tag", delivery.DeliveryTag)

		// log.Println("delivery.Type: ", delivery.RoutingKey)
		switch delivery.RoutingKey {
		case "payment.seccussed":
			c.updateStatus(ctx, delivery, delivery.RoutingKey)
		case "payment.failed":
			c.updateStatus(ctx, delivery, delivery.RoutingKey)
		case "inventory.created":
			c.inventoryCreated(ctx, delivery)
		case "inventory.updated":
			c.inventoryUpdated(ctx, delivery)
		case "inventory.notEnough":
			c.updateStatus(ctx, delivery, delivery.RoutingKey)
		case "inventory.reserved":
			c.updateStatus(ctx, delivery, delivery.RoutingKey)
		default:
			c.handleUnknownDelivery(delivery)
		}
	}
}

func (c *appServer) inventoryCreated(ctx context.Context, delivery amqp091.Delivery) {
	var payload constant.Product
	// Extract trace context
	ctx_ := otel.GetTextMapPropagator().Extract(
		ctx,
		rabbitmq.AMQPHeaderCarrier(delivery.Headers),
	)

	err := json.Unmarshal(delivery.Body, &payload)
	if err != nil {
		slog.Error("failed to Unmarshal", err)
	}
	// tracer := otel.Tracer("order-service")
	// ctx2, span := tracer.Start(ctx_, "Create_Product")
	if err = c.orderService.CreateProduct(ctx_, &payload); err != nil {
		slog.Error("failed to created order_products", err)
		c.rejectDelivery(delivery)
		return
	}
	// span.End()
	c.ackDelivery(delivery)
}

func (c *appServer) updateStatus(ctx context.Context, delivery amqp091.Delivery, routingKey string) {
	var payload int

	ctx_ := otel.GetTextMapPropagator().Extract(
		ctx,
		rabbitmq.AMQPHeaderCarrier(delivery.Headers),
	)

	err := json.Unmarshal(delivery.Body, &payload)
	if err != nil {
		slog.Error("failed to Unmarshal", err)
	}

	var column, status string
	switch routingKey {
	case "payment.seccussed":
		column = "payment_status"
		status = "seccussed"
	case "payment.failed":
		column = "payment_status"
		status = "failed"
	case "inventory.notEnough":
		column = "status"
		status = "failed"
	case "inventory.reserved":
		column = "payment_status"
		status = "pending"
	}

	if err := c.orderService.UpdateStatus(ctx_, payload, column, status); err != nil {
		slog.Error("failed to update order ststus", err)
		c.rejectDelivery(delivery)
		return
	}

	c.ackDelivery(delivery)
}

func (c *appServer) inventoryUpdated(ctx context.Context, delivery amqp091.Delivery) {
	var payload constant.Product
	ctx_ := otel.GetTextMapPropagator().Extract(
		ctx,
		rabbitmq.AMQPHeaderCarrier(delivery.Headers),
	)

	err := json.Unmarshal(delivery.Body, &payload)
	if err != nil {
		slog.Error("failed to Unmarshal", err)
	}

	if err := c.orderService.UpdateProduct(ctx_, &payload); err != nil {
		slog.Error("failed to update order_products", err)
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
	}
}
