package app

import (
	"cart/v1/internal/constant"
	"cart/v1/internal/service"
	"context"
	"encoding/json"
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
	cartService service.CartService
}

func NewWorker(cartService service.CartService) AppServer {
	return &appServer{
		cartService: cartService,
	}
}

func (c *appServer) Worker(ctx context.Context, messages <-chan amqp091.Delivery) {
	for delivery := range messages {

		log.Println("delivery.Type: ", delivery.RoutingKey)
		switch delivery.RoutingKey {
		case "inventory.reserved.buynow":
			c.DeleteCart(ctx, delivery)
		default:
			c.handleUnknownDelivery(delivery)
		}
	}
}

func (c *appServer) DeleteCart(ctx context.Context, delivery amqp091.Delivery) {
	var payload constant.PaymentData
	ctx_ := otel.GetTextMapPropagator().Extract(
		ctx,
		rabbitmq.AMQPHeaderCarrier(delivery.Headers),
	)

	err := json.Unmarshal(delivery.Body, &payload)
	if err != nil {
		slog.Error("failed to Unmarshal", err)
	}

	if payload.OrderSource == "cart" {
		if err := c.cartService.DeleteCart(ctx_, payload.UserID); err != nil {
			slog.Error("failed to delete cart", err)
			c.rejectDelivery(delivery)
			return
		}
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
