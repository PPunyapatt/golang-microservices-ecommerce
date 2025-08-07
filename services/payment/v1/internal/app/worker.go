package app

import (
	"context"
	"encoding/json"
	"log/slog"
	"payment/v1/internal/constant"
	"payment/v1/internal/service"

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
		slog.Info("processDeliveries", "delivery_tag", delivery.DeliveryTag)
		slog.Info("received", "delivery_type", delivery.Type)

		switch delivery.Type {
		case "payment.process":
			var payload constant.PaymentRequest
			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				slog.Error("failed to Unmarshal", err)
			}

			if err = c.paymentService.ProcessPayment(ctx, payload.OrderID, payload.AmountPrice, payload.UserID); err != nil {
				if err = delivery.Reject(false); err != nil {
					slog.Error("failed to delivery.Reject", err)
				}

				slog.Error("failed to process delivery", err)
			} else {
				err = delivery.Ack(false)
				if err != nil {
					slog.Error("failed to acknowledge delivery", err)
				}
			}
		}
	}
}
