package app

import (
	"context"
	"encoding/json"
	"log"
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

		log.Println("delivery.Type: ", delivery.RoutingKey)
		switch delivery.RoutingKey {
		case "payment.seccussed":
			c.paymentSeccussed(ctx, delivery)
		case "inventory.created":
			c.inventoryCreated(ctx, delivery)
		case "inventory.updated":
			c.inventoryUpdated(ctx, delivery)
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

func (c *appServer) paymentSeccussed(ctx context.Context, delivery amqp091.Delivery) {
	var payload constant.UpdateStatus

	err := json.Unmarshal(delivery.Body, &payload)
	if err != nil {
		slog.Error("failed to Unmarshal", err)
	}

	if err := c.orderService.UpdateStatus(ctx); err != nil {
		slog.Error("failed to update order ststus", err)
		c.rejectDelivery(delivery)
		return
	}

	c.ackDelivery(delivery)
}

func (c *appServer) inventoryUpdated(ctx context.Context, delivery amqp091.Delivery) {
	var payload constant.Product

	err := json.Unmarshal(delivery.Body, &payload)
	if err != nil {
		slog.Error("failed to Unmarshal", err)
	}

	if err := c.orderService.UpdateProduct(ctx, &payload); err != nil {
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
