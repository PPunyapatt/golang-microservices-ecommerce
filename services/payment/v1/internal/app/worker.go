package app

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"payment/v1/internal/constant"
	"payment/v1/internal/service"
	"regexp"

	"github.com/rabbitmq/amqp091-go"
)

type AppServer interface {
	Worker(ctx context.Context, messages <-chan amqp091.Delivery)
}

type appServer struct {
	paymentService service.PaymentService
}

func NewWorker(paymentService service.PaymentService) AppServer {
	return &appServer{
		paymentService: paymentService,
	}
}

func (c *appServer) Worker(ctx context.Context, messages <-chan amqp091.Delivery) {
	for delivery := range messages {

		log.Println("delivery.Type: ", delivery.RoutingKey)
		switch delivery.RoutingKey {
		case "inventory.reserved.buynow":
			c.ProcessPayment(ctx, delivery, delivery.RoutingKey)
		case "inventory.reserved.cart":
			c.ProcessPayment(ctx, delivery, delivery.RoutingKey)
		default:
			c.handleUnknownDelivery(delivery)
		}
	}
}

func (c *appServer) ProcessPayment(ctx context.Context, delivery amqp091.Delivery, routingKey string) {
	var payload constant.PaymentRequest
	err := json.Unmarshal(delivery.Body, &payload)
	if err != nil {
		slog.Error("failed to Unmarshal", err)
	}

	re := regexp.MustCompile(`[^.]+$`)
	key := re.FindString(routingKey)
	payload.OrderSource = key

	// if err = c.paymentService.ProcessPayment(ctx, int32(payload.OrderID), payload.TotalPrice); err != nil {
	if err = c.paymentService.ProcessPayment(ctx, &payload); err != nil {
		slog.Error("failed to process payment", err)
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
