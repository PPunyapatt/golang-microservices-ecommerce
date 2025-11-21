package app

import (
	"package/rabbitmq/constant"
	"payment/v1/internal/service"
)

func PaymentConsumer(paymentService service.PaymentService) []constant.ConsumerConfig {
	app := NewWorker(paymentService)
	configs := []constant.ConsumerConfig{
		{
			QueueName: "payment.queue",
			Bindings: []*constant.Queue{
				{
					Exchange: "inventory.exchange",
					Routing:  "inventory.reserved",
				},
				{
					Exchange: "order.dlx",
					Routing:  "order.timeout",
				},
			},
			DeadLetter: &constant.Queue{
				Exchange: "payment.dlx",
				Routing:  "payment.failed",
			},
			WorkerPoolSize: 1,
			Handler:        app.Worker,
			StartWorker:    true,
		},
	}
	return configs
}
