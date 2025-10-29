package rabbitmq

import (
	"context"
	"log/slog"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	_retryTimes     = 5
	_backOffSeconds = 2
)

type RabbitMQ struct {
	Conn *amqp.Connection
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
