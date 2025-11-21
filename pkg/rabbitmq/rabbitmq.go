package rabbitmq

import (
	"context"
	"log"
	"log/slog"
	"package/rabbitmq/constant"
	"package/rabbitmq/consumer"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	_retryTimes     = 5
	_backOffSeconds = 2
)

type RabbitMQ struct {
	Conn *amqp.Connection
}

type consumerManager struct {
	conn      *amqp091.Connection
	cancelCtx context.CancelFunc
	configs   []constant.ConsumerConfig
}

func NewRabbitMQConnection(ctx context.Context, rabbitMqURL string) (*RabbitMQ, error) {
	var (
		counts int64
	)

	for {
		select {
		case <-ctx.Done():
			slog.Warn("🛑 connection canceled by context")
			return nil, ctx.Err()
		default:
			// connection, err := amqp.Dial(string(rabbitMqURL))
			start := time.Now()
			connection, err := amqp.DialConfig(string(rabbitMqURL), amqp.Config{
				Heartbeat: 10 * time.Second,
			})

			if err != nil {
				slog.Error("failed to connect to RabbitMq...", err.Error(), rabbitMqURL)
				counts++
			} else {
				slog.Info("📫 connected to rabbitmq 🎉")
				return &RabbitMQ{Conn: connection}, nil
			}

			duration_connect := time.Since(start).Seconds()

			if counts > _retryTimes {
				slog.Error("failed to retry", "error", err.Error())
				return nil, err
			}

			slog.Info("Backing off for 2 seconds...")
			time.Sleep(_backOffSeconds * time.Second)
			duration_full := time.Since(start).Seconds()

			slog.Debug("Time duration",
				"duration_connect", duration_connect,
				"duration_full", duration_full,
				"count", counts,
			)
		}
	}
}

func (r *RabbitMQ) HandleGracefulShutdown(ctx context.Context, wg *sync.WaitGroup) {
	<-ctx.Done()
	r.Conn.Close()
	slog.Info("🛑 shutting down rabbitmq connection...")
	wg.Done()
}

//----------------------------------------------------------------------------------------------------------------------------
//----------------------------------------------------------------------------------------------------------------------------

func NewConsumerManager(conn *amqp091.Connection, configs []constant.ConsumerConfig) *consumerManager {
	return &consumerManager{
		conn:    conn,
		configs: configs,
	}
}

func (c *consumerManager) InitConsumers(ctx context.Context, conn *amqp091.Connection) {
	if c.cancelCtx != nil {
		c.cancelCtx()
	}

	ctxCancel, cancel := context.WithCancel(ctx)
	c.cancelCtx = cancel

	for _, config := range c.configs {
		cons := consumer.NewConsumer(conn)
		cons.Configure(
			consumer.QueueName(config.QueueName),
			consumer.QueueProperties(config.Bindings),
			consumer.TopicType("topic"),
		)

		if config.DeadLetter != nil {
			cons.Configure(consumer.QueueDeadLetter(config.DeadLetter))
		}

		if config.WorkerPoolSize > 0 {
			cons.Configure(consumer.WorkerPoolSize(config.WorkerPoolSize))
		}

		if config.StartWorker {
			go cons.StartConsumer(ctxCancel, config.Handler)
		} else {
			go cons.StartConsumer(ctxCancel, nil)
		}
	}
}

func (c *consumerManager) InitConsumerWithReconnection(ctx context.Context, RabbitMQUrl string) {
	c.InitConsumers(ctx, c.conn)

	go func() {
		backoff := time.Second
		maxBackoff := 30 * time.Second

		errCh := c.conn.NotifyClose(make(chan *amqp091.Error))
		for {
			select {
			case <-ctx.Done():
				if c.cancelCtx != nil {
					c.cancelCtx()
				}
				slog.Info("RabbitMQ connection closed: context canceled")
				return
			case err := <-errCh:
				log.Printf("connection closed: %+v", errors.WithStack(err))
				if c.cancelCtx != nil {
					c.cancelCtx()
				}

				for {
					rb, err := NewRabbitMQConnection(ctx, RabbitMQUrl)
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
					c.InitConsumers(ctx, c.conn)
					log.Println(" ----------- Reconnect successed ----------- ")
					break
				}
			}
		}
	}()

}
