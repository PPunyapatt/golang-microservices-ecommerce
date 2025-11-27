package app

import (
	"context"
	"encoding/json"
	"inventories/v1/internal/constant"
	"inventories/v1/internal/services"
	"log"
	"log/slog"
	"package/rabbitmq"
	"regexp"
	"sync"

	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

type AppServer interface {
	Worker(ctx context.Context, message amqp091.Delivery)
}

type appServer struct {
	inventoryService services.InventoryServie
}

var orderPool = sync.Pool{
	New: func() interface{} {
		return new(constant.Order)
	},
}

func NewWorker(inventoryService services.InventoryServie) AppServer {
	return &appServer{
		inventoryService: inventoryService,
	}
}

func (c *appServer) Worker(ctx context.Context, message amqp091.Delivery) {
	log.Println("delivery.Type: ", message.RoutingKey)
	reserveStockKeys := map[string]struct{}{
		"order.created": {},
	}

	cutOrReleaseStockKeys := map[string]struct{}{
		"order.payment.successed": {},
		"order.payment.failed":    {},
		"order.timeout":           {},
	}

	if _, ok := reserveStockKeys[message.RoutingKey]; ok {
		c.ReserveStock(ctx, message, message.RoutingKey)
	} else if _, ok := cutOrReleaseStockKeys[message.RoutingKey]; ok {
		c.CutOrReleaseStock(ctx, message, message.RoutingKey)
	} else {
		c.handleUnknownDelivery(message)
	}
}

func (c *appServer) ReserveStock(ctx context.Context, delivery amqp091.Delivery, routingKey string) {
	ctx_ := otel.GetTextMapPropagator().Extract(
		ctx,
		rabbitmq.AMQPHeaderCarrier(delivery.Headers),
	)

	// var payload constant.Order
	payload := orderPool.Get().(*constant.Order)
	defer orderPool.Put(payload)

	err := json.Unmarshal(delivery.Body, &payload)
	if err != nil {
		slog.Error("failed to Unmarshal", "err", err.Error())
	}

	re := regexp.MustCompile(`[^.]+$`)
	key := re.FindString(routingKey)

	if err = c.inventoryService.ReserveStock(ctx_, payload, key); err != nil {
		slog.Error("failed to reserve stock", "err", err.Error())
		c.rejectDelivery(delivery)
		return
	}

	c.ackDelivery(delivery)
}

func (c *appServer) CutOrReleaseStock(ctx context.Context, delivery amqp091.Delivery, routingKey string) {
	var payload []*constant.Item
	err := json.Unmarshal(delivery.Body, &payload)
	if err != nil {
		slog.Error("failed to Unmarshal", "err", err.Error())
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
