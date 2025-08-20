package app

import (
	"package/rabbitmq/constant"
	"package/rabbitmq/consumer"
	"payment/v1/internal/service"

	"github.com/rabbitmq/amqp091-go"
)

func InitConsumers(paymentService service.PaymentService, conn *amqp091.Connection) {
	paymentQueue := []*constant.Queue{
		{
			Exchange: "inventory.exchange",
			Routing:  "inventory.reserved",
		},
	}

	paymentDLQueue := &constant.Queue{
		Exchange: "payment.dlx",
		Routing:  "payment.failed",
	}

	paymentConsumer := consumer.NewConsumer(conn)
	paymentConsumer.Configure(
		consumer.QueueName("payment.queue"),
		consumer.QueueProperties(paymentQueue),
		consumer.QueueDeadLetter(paymentDLQueue),
		consumer.WorkerPoolSize(1),
		consumer.TopicType("topic"),
	)

	app := NewWorker(paymentService)

	go paymentConsumer.StartConsumer(app.Worker)
}
