package publisher

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	_publishMandatory = false
	_publishImmediate = false

	_exchangeName    = "orders-exchange"
	_bindingKey      = "orders-routing-key"
	_messageTypeName = "ordered"
)

type EventPublisher interface {
	Configure(...Option) EventPublisher
	Publish(context.Context, []byte, string) error
}

type Publisher struct {
	exchangeName, bindingKey string
	messageTypeName          string
	amqpChan                 *amqp.Channel
	amqpConn                 *amqp.Connection
}

func NewPublisher(amqpConn *amqp.Connection) EventPublisher {
	ch, err := amqpConn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	return &Publisher{
		exchangeName: _exchangeName,
		bindingKey:   _bindingKey,
		amqpConn:     amqpConn,
		amqpChan:     ch,
	}
}

func (p *Publisher) Configure(opts ...Option) EventPublisher {
	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Publisher) Publish(ctx context.Context, body []byte, contentType string) error {
	return nil
}
