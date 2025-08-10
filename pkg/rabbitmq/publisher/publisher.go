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
	exchangeName, routingKeys []string
	topicType                 string
	amqpConn                  *amqp.Connection
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

func (p *Publisher) Publish(ctx context.Context, body []byte, routingKey string) error {
	ch, err := p.amqpConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	for _, exchangeName := range p.exchangeName {
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

		for _, routingName := range p.routingKeys {
			if err := ch.PublishWithContext(
				ctx,
				exchangeName,
				routingName,
				false, // mandatory
				false, // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        body,
				},
			); err != nil {
				return err
			}
		}
	}

	return nil
}
