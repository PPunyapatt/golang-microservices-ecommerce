package app

import (
	"cart/v1/internal/service"
	"context"
	"log"
	"os"
	"package/rabbitmq"
	"package/rabbitmq/constant"
	"package/rabbitmq/consumer"
	"time"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
)

type ConsumerManager struct {
	conn        *amqp091.Connection
	cancelCtx   context.CancelFunc
	cartService service.CartService
}

func NewInitConsumer(cartService service.CartService, conn *amqp091.Connection) *ConsumerManager {
	return &ConsumerManager{
		cartService: cartService,
		conn:        conn,
	}
}

func (c *ConsumerManager) InitConsumer(conn *amqp091.Connection) {
	if c.cancelCtx != nil {
		c.cancelCtx()
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.cancelCtx = cancel

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
		consumer.WorkerPoolSize(3),
		consumer.TopicType("topic"),
	)

	app := NewWorker(c.cartService)
	go cartConsumer.StartConsumer(ctx, app.Worker)

}

func (c *ConsumerManager) InitConsumerWithReconnection(ctx context.Context) {
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
				rb, err := rabbitmq.NewRabbitMQConnection(ctx, url)
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
				c.conn = rb.Conn
				c.InitConsumer(c.conn)
				log.Println(" ----------- Reconnect successed ----------- ")
				break
			}

		}
	}()

}
