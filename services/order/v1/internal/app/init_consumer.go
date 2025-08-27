package app

import (
	"context"
	"log"
	"order/v1/internal/service"
	"os"
	"package/rabbitmq"
	"package/rabbitmq/constant"
	"package/rabbitmq/consumer"
	"time"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
)

type ConsumerManager struct {
	conn         *amqp091.Connection
	cancelCtx    context.CancelFunc
	orderService service.OrderService
}

func NewInitConsumer(orderService service.OrderService, conn *amqp091.Connection) *ConsumerManager {
	return &ConsumerManager{
		orderService: orderService,
		conn:         conn,
	}
}

func (c *ConsumerManager) InitConsumer(conn *amqp091.Connection) {

	if c.cancelCtx != nil {
		c.cancelCtx()
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.cancelCtx = cancel

	// ---------------- Order Queue ----------------
	// ---------------------------------------------
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

	orderConsumer := consumer.NewConsumer(c.conn)
	orderConsumer.Configure(
		consumer.QueueName("order.queue"),
		consumer.QueueProperties(orderQueues),
		consumer.WorkerPoolSize(1),
		consumer.TopicType("topic"),
	)

	// ---------------- Order Dead Letter Queue ----------------
	// ---------------------------------------------------------
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
	orderDLconsumer := consumer.NewConsumer(c.conn)
	orderDLconsumer.Configure(
		consumer.QueueName("order.dlq"),
		consumer.QueueProperties(orderDLQueues),
		consumer.WorkerPoolSize(1),
		consumer.TopicType("topic"),
	)

	// ---------------- Order Delay Queue ----------------
	// ---------------------------------------------------
	orderBind := []*constant.Queue{
		{
			Exchange: "inventory.exchange",
			Routing:  "inventory.reserved",
		},
	}

	DeadLetterDestination := &constant.Queue{
		Exchange: "order.dlx",
		Routing:  "order.timeout",
	}

	orderDelayconsumer := consumer.NewConsumer(c.conn)
	orderDelayconsumer.Configure(
		consumer.QueueName("order.delay.queue"),
		consumer.QueueProperties(orderBind),
		consumer.QueueDeadLetter(DeadLetterDestination),
		consumer.TopicType("topic"),
	)

	app_ := NewWorker(c.orderService)
	go orderConsumer.StartConsumer(ctx, app_.Worker)

	appDlx := NewWorker(c.orderService)
	go orderDLconsumer.StartConsumer(ctx, appDlx.Worker)

	go orderDelayconsumer.StartConsumer(ctx, nil)

}

func (c *ConsumerManager) InitConsumerWithReconnection() {
	c.InitConsumer(c.conn)

	go func() {
		backoff := time.Second
		maxBackoff := 30 * time.Second
		for {
			err := <-c.conn.NotifyClose(make(chan *amqp091.Error))
			if err != nil {
				log.Printf("connection closed: %+v", errors.WithStack(err))
			}

			for {
				url := os.Getenv("RABBITMQ")
				newConn, err := rabbitmq.NewRabbitMQConnection(url)
				if err != nil {
					log.Printf("reconnect failed: %+v", errors.WithStack(err))
					log.Printf("retry in %s ...", backoff)
					time.Sleep(backoff)

					backoff *= 2
					if backoff > maxBackoff {
						backoff = maxBackoff
					}
					continue
				}

				backoff = time.Second
				c.conn = newConn
				c.InitConsumer(c.conn)
				log.Println(" ----------- Reconnect successed ----------- ")
				break
			}

		}
	}()

}
