package app

import (
	"order/v1/internal/service"
	"package/rabbitmq/constant"
	"package/rabbitmq/consumer"

	"github.com/rabbitmq/amqp091-go"
)

func InitConsumer(orderService service.OrderService, conn *amqp091.Connection) {
	orderQueues := []*constant.Queue{
		{
			Exchange: "inventory.exchange",
			Routing:  "inventory.*",
		},
		{
			Exchange: "payment.exchange",
			Routing:  "payment.*",
		},
	}

	orderConsumer := consumer.NewConsumer(conn)
	orderConsumer.Configure(
		consumer.QueueName("order.queue"),
		consumer.QueueProperties(orderQueues),
		consumer.WorkerPoolSize(1),
		consumer.TopicType("topic"),
	)

	orderDLQueues := []*constant.Queue{
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
	}
	orderDLconsumer := consumer.NewConsumer(conn)
	orderDLconsumer.Configure(
		consumer.QueueName("order.dlq"),
		consumer.QueueProperties(orderDLQueues),
		consumer.WorkerPoolSize(1),
		consumer.TopicType("topic"),
	)

	orderDelay := []*constant.Queue{
		{
			Exchange: "order.exchange",
			Routing:  "order.created",
		},
	}

	orderDelayDL := &constant.Queue{
		Exchange: "order.dlx",
		Routing:  "order.time_out_15_min",
	}

	orderDelayconsumer := consumer.NewConsumer(conn)
	orderDelayconsumer.Configure(
		consumer.QueueName("order.delay.queue"),
		consumer.QueueProperties(orderDelay),
		consumer.QueueDeadLetter(orderDelayDL),
		consumer.WorkerPoolSize(1),
		consumer.TopicType("topic"),
	)

	app_ := NewWorker(orderService)
	go orderConsumer.StartConsumer(app_.Worker)

	appDlx := NewWorker(orderService)
	go orderDLconsumer.StartConsumer(appDlx.Worker)

}
