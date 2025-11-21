package app

import (
	"order/v1/internal/service"
	"package/rabbitmq/constant"
)

func OrderConsumer(orderService service.OrderService) []constant.ConsumerConfig {
	app := NewWorker(orderService)
	appDlx := NewWorker(orderService)
	configs := []constant.ConsumerConfig{
		{
			QueueName: "order.queue",
			Bindings: []*constant.Queue{
				{
					Exchange: "inventory.exchange",
					Routing:  "inventory.*",
				},
				{
					Exchange: "payment.exchange",
					Routing:  "payment.*",
				},
			},
			WorkerPoolSize: 3,
			Handler:        app.Worker,
			StartWorker:    true,
		},
		{
			QueueName: "order.dlq",
			Bindings: []*constant.Queue{
				{
					Exchange: "inventory.dlx",
					Routing:  "inventory.failed",
				},
				{
					Exchange: "payment.dlx",
					Routing:  "payment.failed",
				},
				{
					Exchange: "order.dlx",
					Routing:  "order.timeout",
				},
			},
			WorkerPoolSize: 3,
			Handler:        appDlx.Worker,
			StartWorker:    true,
		},
		{
			QueueName: "order.delay.queue",
			Bindings: []*constant.Queue{
				{
					Exchange: "order.dlx",
					Routing:  "order.timeout",
				},
			},
			DeadLetter: &constant.Queue{
				Exchange: "order.dlx",
				Routing:  "order.timeout",
			},
			WorkerPoolSize: 3,
			StartWorker:    false,
		},
	}
	return configs
}
