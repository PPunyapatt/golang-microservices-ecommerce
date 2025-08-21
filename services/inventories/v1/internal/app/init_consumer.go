package app

import (
	"inventories/v1/internal/services"
	"package/rabbitmq/constant"
	"package/rabbitmq/consumer"

	"github.com/rabbitmq/amqp091-go"
)

func InitConsumer(inventoryService services.InventoryServie, conn *amqp091.Connection) {
	inventoryQueues := []*constant.Queue{
		{
			Exchange: "order.exchange",
			Routing:  "order.#",
		},
	}

	inventoryDLQueues := &constant.Queue{
		Exchange: "inventory.exchange",
		Routing:  "inventory.failed",
	}

	inventoryConsumer := consumer.NewConsumer(conn, true)
	inventoryConsumer.Configure(
		consumer.QueueProperties(inventoryQueues),
		consumer.QueueDeadLetter(inventoryDLQueues),
		consumer.QueueName("inventory.queue"),
		consumer.WorkerPoolSize(1),
		consumer.TopicType("topic"),
	)

	app := NewWorker(inventoryService)
	go inventoryConsumer.StartConsumer(app.Worker)
}
