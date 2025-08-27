package app

import (
	"context"
	"inventories/v1/internal/services"
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
	conn             *amqp091.Connection
	cancelCtx        context.CancelFunc
	inventoryService services.InventoryServie
}

func NewInitConsumer(inventoryService services.InventoryServie, conn *amqp091.Connection) *ConsumerManager {
	return &ConsumerManager{
		inventoryService: inventoryService,
		conn:             conn,
	}
}

func (c *ConsumerManager) InitConsumer(conn *amqp091.Connection) {
	if c.cancelCtx != nil {
		c.cancelCtx()
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.cancelCtx = cancel

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

	app := NewWorker(c.inventoryService)
	go inventoryConsumer.StartConsumer(ctx, app.Worker)
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
