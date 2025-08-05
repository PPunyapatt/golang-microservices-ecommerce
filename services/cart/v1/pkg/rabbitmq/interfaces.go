package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type IRabbitMQ interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool) (<-chan amqp.Delivery, error)
}
