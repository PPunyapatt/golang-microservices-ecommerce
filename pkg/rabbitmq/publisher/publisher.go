package publisher

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

type EventPublisher interface {
	Configure(...Option) EventPublisher
	Publish(context.Context, []byte, string, string, amqp091.Table, ...int) error
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

func (p *Publisher) Publish(ctx context.Context, body []byte, exchangeName, routingKey string, headers amqp091.Table, ttl ...int) error {
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

	amqpPublishing := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
		Headers:     headers,
	}

	if len(ttl) > 0 {
		timeExpire := time.Duration(ttl[0]) * time.Minute
		expire := strconv.FormatInt(timeExpire.Milliseconds(), 10)
		amqpPublishing.Expiration = expire
	}

	if err := p.ch.PublishWithContext(
		ctx,
		exchangeName,
		routingKey,
		false, // mandatory
		false, // immediate
		amqpPublishing,
	); err != nil {
		return err
	}

	return nil
}
