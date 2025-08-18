package app

import (
	"context"
	"encoding/json"
	"log/slog"
	"order/v1/internal/service"
	"package/rabbitmq"

	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

type AppServerDeadLetter interface {
	Worker(ctx context.Context, messages <-chan amqp091.Delivery)
}

type appServerDeadLetter struct {
	orderService service.OrderService
}

func NewWorkerDeadLetter(orderService service.OrderService) AppServerDeadLetter {
	return &appServerDeadLetter{
		orderService: orderService,
	}
}

func (c *appServerDeadLetter) Worker(ctx context.Context, messages <-chan amqp091.Delivery) {
	for delivery := range messages {
		slog.Info("processDeliveries", "delivery_tag", delivery.DeliveryTag)

		// log.Println("delivery.Type: ", delivery.RoutingKey)
		switch delivery.RoutingKey {
		case "payment.failed":
			c.updateStatus(ctx, delivery, delivery.RoutingKey)

		case "inventory.notEnough":
			c.updateStatus(ctx, delivery, delivery.RoutingKey)

		default:
			c.handleUnknownDelivery(delivery)
		}
	}
}

func (c *appServerDeadLetter) updateStatus(ctx context.Context, delivery amqp091.Delivery, routingKey string) {
	var orderID int

	ctx_ := otel.GetTextMapPropagator().Extract(
		ctx,
		rabbitmq.AMQPHeaderCarrier(delivery.Headers),
	)

	err := json.Unmarshal(delivery.Body, &orderID)
	if err != nil {
		slog.Error("failed to Unmarshal", err)
	}

	type updateRule struct {
		updates      map[string]interface{}
		publishStock bool
	}

	rules := map[string]updateRule{
		"payment.failed": updateRule{
			updates: map[string]interface{}{
				"status":         "failed",
				"payment_status": "payment_failed",
			},
			publishStock: true,
		},
		"inventory.failed": updateRule{
			updates: map[string]interface{}{
				"status": "failed",
			},
		},
	}

	if err := c.orderService.UpdateStatus(ctx_, orderID, rules[routingKey].updates); err != nil {
		slog.Error("failed to update order ststus", err)
		c.rejectDelivery(delivery)
		return
	}

	if rules[routingKey].publishStock {
		if err := c.orderService.PushEventCutorReleaseStock(ctx_, orderID, routingKey); err != nil {
			slog.Error("failed to update CutorReleaseStockEvent", err)
			c.rejectDelivery(delivery)
			return
		}
	}

	c.ackDelivery(delivery)
}

// -------------------------- Handler Error --------------------------
func (c *appServerDeadLetter) handleUnknownDelivery(delivery amqp091.Delivery) {
	slog.Warn("unknown delivery routing key", "key", delivery.RoutingKey)
	c.rejectDelivery(delivery)
}

func (c *appServerDeadLetter) rejectDelivery(delivery amqp091.Delivery) {
	if err := delivery.Reject(false); err != nil {
		slog.Error("failed to delivery.Reject", err)
	}
}

func (c *appServerDeadLetter) ackDelivery(delivery amqp091.Delivery) {
	err := delivery.Ack(false)
	if err != nil {
		slog.Error("failed to acknowledge delivery", err)
	}
}
