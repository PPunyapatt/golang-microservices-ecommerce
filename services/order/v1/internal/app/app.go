package app

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"order/v1/internal/constant"
	"order/v1/internal/service"
	"package/rabbitmq"
	"time"

	"github.com/goforj/godump"
	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

type AppServer interface {
	Worker(ctx context.Context, message amqp091.Delivery)
}

type appServer struct {
	orderService service.OrderService
}

func NewWorker(orderService service.OrderService) AppServer {
	return &appServer{
		orderService: orderService,
	}
}

type orderPayload struct {
	OrderID     int    `json:"order_id"`
	OrderSource string `json:"order_source"`
}

func (c *appServer) Worker(ctx context.Context, message amqp091.Delivery) {
	// for delivery := range messages {
	slog.Info("processDeliveries", "delivery_tag", message.DeliveryTag)

	updateStatusKey := map[string]struct{}{
		"payment.successed":  {},
		"inventory.reserved": {},
		"payment.failed":     {},
		"inventory.failed":   {},
	}

	log.Println("message.Type: ", message.RoutingKey)
	switch message.RoutingKey {
	case "inventory.event":
		c.inventoryEvent(ctx, message)

	case "order.timeout":
		c.checkAndUpdateStatus(ctx, message, message.RoutingKey)
	default:
		if _, ok := updateStatusKey[message.RoutingKey]; ok {
			c.updateStatus(ctx, message, message.RoutingKey)
		} else {
			c.handleUnknownMessage(message)
		}
	}
	// }
}

func (c *appServer) updateStatus(ctx context.Context, delivery amqp091.Delivery, routingKey string) {
	var order orderPayload

	ctx_ := otel.GetTextMapPropagator().Extract(
		ctx,
		rabbitmq.AMQPHeaderCarrier(delivery.Headers),
	)

	err := json.Unmarshal(delivery.Body, &order)
	if err != nil {
		slog.Error("failed to Unmarshal", err)
	}

	type updateRule struct {
		updates      map[string]interface{}
		publishStock bool
	}

	rules := map[string]updateRule{
		"payment.successed": updateRule{
			updates: map[string]interface{}{
				"updated_at":     time.Now().UTC(),
				"status":         "successed",
				"payment_status": "payment_successed",
			},
			publishStock: true,
		},
		"inventory.reserved": updateRule{
			updates: map[string]interface{}{
				"updated_at":     time.Now().UTC(),
				"status":         "reserved",
				"payment_status": "pending",
			},
		},
		"payment.failed": updateRule{
			updates: map[string]interface{}{
				"updated_at":     time.Now().UTC(),
				"status":         "failed",
				"payment_status": "payment_failed",
			},
			publishStock: true,
		},
		"inventory.failed": updateRule{
			updates: map[string]interface{}{
				"updated_at": time.Now().UTC(),
				"status":     "failed",
			},
		},
	}

	if err := c.orderService.UpdateStatus(ctx_, order.OrderID, rules[routingKey].updates); err != nil {
		slog.Error("failed to update order ststus", err)
		c.rejectDelivery(delivery)
		return
	}

	if rules[routingKey].publishStock {
		if err := c.orderService.PushEventCutorReleaseStock(ctx_, order.OrderID, routingKey); err != nil {
			slog.Error("failed to update order_products", err)
			c.rejectDelivery(delivery)
			return
		}
	}

	c.ackDelivery(delivery)
}

func (c *appServer) inventoryEvent(ctx context.Context, message amqp091.Delivery) {
	var payload map[string]constant.OrderProduct
	err := json.Unmarshal(message.Body, &payload)
	if err != nil {
		slog.Error("failed to Unmarshal", err)
		c.rejectDelivery(message)
		return
	}

	data := payload["payload"]
	godump.Dump(data)
	if err := c.orderService.ProductUpdate(ctx, &data); err != nil {
		slog.Error("failed to update order_product", err)
		c.rejectDelivery(message)
		return
	}

	c.ackDelivery(message)
}

func (c *appServer) checkAndUpdateStatus(ctx context.Context, delivery amqp091.Delivery, routingKey string) {
	var order orderPayload
	ctx_ := otel.GetTextMapPropagator().Extract(
		ctx,
		rabbitmq.AMQPHeaderCarrier(delivery.Headers),
	)

	err := json.Unmarshal(delivery.Body, &order)
	if err != nil {
		slog.Error("failed to Unmarshal", err)
	}

	if err := c.orderService.PushEventCutorReleaseStock(ctx_, order.OrderID, routingKey); err != nil {
		slog.Error("failed to publish event", err)
		c.rejectDelivery(delivery)
		return
	}

	err = c.orderService.CheckAndUpdateStatus(ctx_, order.OrderID)
	if err != nil {
		slog.Error("failed to update order status", err)
		c.rejectDelivery(delivery)
		return
	}

	c.ackDelivery(delivery)
}

// -------------------------- Handler Error --------------------------
func (c *appServer) handleUnknownMessage(delivery amqp091.Delivery) {
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
