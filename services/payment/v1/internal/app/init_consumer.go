package app

import (
	"context"
	"log"
	"os"
	"package/rabbitmq"
	"package/rabbitmq/constant"
	"package/rabbitmq/consumer"
	"payment/v1/internal/service"
	"time"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
)

type ConsumerManager struct {
	conn           *amqp091.Connection
	cancelCtx      context.CancelFunc
	paymentService service.PaymentService
}

func NewInitConsumer(paymentService service.PaymentService, conn *amqp091.Connection) *ConsumerManager {
	return &ConsumerManager{
		paymentService: paymentService,
		conn:           conn,
	}
}

func (c *ConsumerManager) InitConsumer(conn *amqp091.Connection) {
	if c.cancelCtx != nil {
		c.cancelCtx()
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.cancelCtx = cancel

	paymentQueue := []*constant.Queue{
		{
			Exchange: "inventory.exchange",
			Routing:  "inventory.reserved",
		},
		{
			Exchange: "order.dlx",
			Routing:  "order.timeout",
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

	app := NewWorker(c.paymentService)

	go paymentConsumer.StartConsumer(ctx, app.Worker)
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
