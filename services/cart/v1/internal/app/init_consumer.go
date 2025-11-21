package app

import (
	"cart/v1/internal/service"
	"package/rabbitmq/constant"
)

func CartConsumer(cartService service.CartService) []constant.ConsumerConfig {
	app := NewWorker(cartService)
	configs := []constant.ConsumerConfig{
		{
			QueueName: "cart.queue",
			Bindings: []*constant.Queue{
				{
					Exchange: "payment.exchange",
					Routing:  "payment.successed",
				},
			},
			WorkerPoolSize: 3,
			Handler:        app.Worker,
			StartWorker:    true,
		},
	}
	return configs
}
