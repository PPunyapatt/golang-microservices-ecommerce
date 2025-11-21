package app

import (
	"inventories/v1/internal/services"
	"package/rabbitmq/constant"
)

func InventoryConsumer(inventoryService services.InventoryServie) []constant.ConsumerConfig {
	app := NewWorker(inventoryService)
	configs := []constant.ConsumerConfig{
		{
			QueueName: "inventory.queue",
			Bindings: []*constant.Queue{
				{
					Exchange: "order.exchange",
					Routing:  "order.#",
				},
			},
			DeadLetter: &constant.Queue{
				Exchange: "inventory.exchange",
				Routing:  "inventory.failed",
			},
			WorkerPoolSize: 1,
			Handler:        app.Worker,
			StartWorker:    true,
		},
	}
	return configs
}
