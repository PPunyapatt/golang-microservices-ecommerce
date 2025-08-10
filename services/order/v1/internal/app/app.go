package app

import (
	"context"
	"encoding/json"
	"log/slog"
	"order/v1/internal/constant"
	"order/v1/proto/order"

	"github.com/rabbitmq/amqp091-go"
)

type AppServer interface {
	Worker(ctx context.Context, messages <-chan amqp091.Delivery)
}

type appServer struct {
	orderService order.OrderServiceServer
}

func NewWorker(orderService order.OrderServiceServer) AppServer {
	return &appServer{
		orderService: orderService,
	}
}

func (c *appServer) Worker(ctx context.Context, messages <-chan amqp091.Delivery) {
	for delivery := range messages {
		slog.Info("processDeliveries", "delivery_tag", delivery.DeliveryTag)
		slog.Info("received", "delivery_type", delivery.Type)

		switch delivery.Type {
		case "payment.seccussed":
			var payload constant.UpdateStatus

			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				slog.Error("failed to Unmarshal", err)
			}

			_, err = c.orderService.UpdateStatus(ctx, &order.UpdateStatusRequest{
				PaymentStatus: payload.PaymentStatus,
				OrderStatus:   payload.OrderStatus,
				OrderID:       int32(payload.OrderID),
			})

			if err != nil {
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
