package app

import (
	"cart/v1/internal/service"
	"package/rabbitmq/constant"
	"package/rabbitmq/consumer"

	"github.com/rabbitmq/amqp091-go"
)

func InitConsumer(cartService service.CartService, conn *amqp091.Connection) {
	cartQueueBiding := []*constant.Queue{
		{
			Exchange: "payment.exchange",
			Routing:  "payment.successed",
		},
	}

	cartConsumer := consumer.NewConsumer(conn)
	cartConsumer.Configure(
		consumer.QueueProperties(cartQueueBiding),
		consumer.QueueName("cart.queue"),
		consumer.WorkerPoolSize(1),
		consumer.TopicType("topic"),
	)

	app := NewWorker(cartService)
	go cartConsumer.StartConsumer(app.Worker)

}
