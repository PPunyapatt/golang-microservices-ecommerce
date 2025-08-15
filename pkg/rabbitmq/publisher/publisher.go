package publisher

import (
	"context"
	"fmt"
	"sync"

	"github.com/rabbitmq/amqp091-go"
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
	Publish(context.Context, []byte, string, string, amqp091.Table) error
}

type Publisher struct {
	routingKeys []string
	topicType   string
	amqpConn    *amqp.Connection
	ch          *amqp.Channel
	mu          sync.Mutex
}

func NewPublisher(amqpConn *amqp.Connection) (EventPublisher, error) {
	var mu sync.Mutex
	ch, err := amqpConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}
	return &Publisher{
		amqpConn: amqpConn,
		ch:       ch,
		mu:       mu,
	}, nil
}

func (p *Publisher) Configure(opts ...Option) EventPublisher {
	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Publisher) Publish(ctx context.Context, body []byte, exchangeName, routingKey string, headers amqp091.Table) error {
	// ch, err := p.amqpConn.Channel()
	// if err != nil {
	// 	return err
	// }
	// defer ch.Close()

	p.mu.Lock()
	defer p.mu.Unlock()

	// var ch *amqp.Channel
	if p.ch.IsClosed() {
		ch, err := p.amqpConn.Channel()
		if err != nil {
			return err
		}
		p.ch = ch
	}

	if err := p.ch.ExchangeDeclare(
		exchangeName, // name
		p.topicType,  // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	); err != nil {
		return err
	}

	if err := p.ch.PublishWithContext(
		ctx,
		exchangeName,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Headers:     headers,
		},
	); err != nil {
		return err
	}

	return nil
}
