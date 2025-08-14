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
	Publish(context.Context, []byte, string, string) error
}

type Publisher struct {
	routingKeys []string
	topicType   string
	amqpConn    *amqp.Connection
}

func NewPublisher(amqpConn *amqp.Connection) EventPublisher {
	return &Publisher{
		// exchangeName: _exchangeName,
		// bindingKey:   _bindingKey,
		amqpConn: amqpConn,
	}
}

func (p *Publisher) Configure(opts ...Option) EventPublisher {
	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Publisher) Publish(ctx context.Context, body []byte, exchangeName, routingKey string) error {
	ch, err := p.amqpConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if err = ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	); err != nil {
		return err
	}

	if err := ch.PublishWithContext(
		ctx,
		exchangeName,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	); err != nil {
		return err
	}

	return nil
}
