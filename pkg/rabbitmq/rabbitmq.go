package rabbitmq

import (
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	_retryTimes     = 5
	_backOffSeconds = 2
)

func NewRabbitMQConnection(rabbitMqURL string) (*amqp.Connection, error) {
	var (
		amqpConn *amqp.Connection
		counts   int64
	)

	for {
		// connection, err := amqp.Dial(string(rabbitMqURL))
		connection, err := amqp.DialConfig(string(rabbitMqURL), amqp.Config{
			Heartbeat: 10 * time.Second,
		})

		if err != nil {
			slog.Error("failed to connect to RabbitMq...", err, rabbitMqURL)
			counts++
		} else {
			amqpConn = connection

			break
		}

		if counts > _retryTimes {
			slog.Error("failed to retry", err)

			return nil, err
		}

		slog.Info("Backing off for 2 seconds...")
		time.Sleep(_backOffSeconds * time.Second)

		continue
	}

	slog.Info("📫 connected to rabbitmq 🎉")

	return amqpConn, nil
}
