package rabbitmq

import (
	"errors"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	_retryTimes     = 5
	_backOffSeconds = 2
)

// type IRabbitMQ interface {
// 	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
// 	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool) (<-chan amqp.Delivery, error)
// }

// type RabbitMQ struct {
// 	amqpChan *amqp.Channel
// 	amqpConn *amqp.Connection
// }

// var _ IRabbitMQ = (*RabbitMQ)(nil)

func NewRabbitMQConnection(rabbitMqURL string) (*amqp.Connection, error) {
	var (
		amqpConn *amqp.Connection
		counts   int64
	)

	var ErrCannotConnectRabbitMQ = errors.New("cannot connect to rabbit")

	for {
		connection, err := amqp.Dial(string(rabbitMqURL))
		if err != nil {
			slog.Error("failed to connect to RabbitMq...", err, rabbitMqURL)
			counts++
		} else {
			amqpConn = connection

			break
		}

		if counts > _retryTimes {
			slog.Error("failed to retry", err)

			return nil, ErrCannotConnectRabbitMQ
		}

		slog.Info("Backing off for 2 seconds...")
		time.Sleep(_backOffSeconds * time.Second)

		continue
	}

	slog.Info("📫 connected to rabbitmq 🎉")

	return amqpConn, nil
}

// func (r *RabbitMQ) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
// 	return nil
// }

// func (r *RabbitMQ) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool) (<-chan amqp.Delivery, error) {
// 	deliveries, err := r.amqpChan.Consume(
// 		queue,
// 		consumer,
// 		autoAck,
// 		exclusive,
// 		noLocal,
// 		noWait,
// 		nil,
// 	)
// 	if err != nil {
// 		log.Fatalf("Failed to register a consumer: %s", err)
// 		return nil, err
// 	}

// 	return deliveries, nil
// }
