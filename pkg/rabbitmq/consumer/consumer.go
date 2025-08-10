package consumer

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	_exchangeKind       = "direct"
	_exchangeDurable    = true
	_exchangeAutoDelete = false
	_exchangeInternal   = false
	_exchangeNoWait     = false

	_queueDurable    = true
	_queueAutoDelete = false
	_queueExclusive  = false
	_queueNoWait     = false

	_prefetchCount  = 5
	_prefetchSize   = 0
	_prefetchGlobal = false

	_consumeAutoAck   = false
	_consumeExclusive = false
	_consumeNoLocal   = false
	_consumeNoWait    = false

	_exchangeName   = "orders-exchange"
	_queueName      = "orders-queue"
	_bindingKey     = "orders-routing-key"
	_consumerTag    = "orders-consumer"
	_workerPoolSize = 2
)

type worker func(ctx context.Context, messages <-chan amqp.Delivery)
type EventConsumer interface {
	Configure(...Option) EventConsumer
	StartConsumer(worker) error
}

type consumer struct {
	queueName, bindingKey, consumerName, topicType string
	workerPoolSize                                 int
	RoutingKeys, exchangeName                      []string
	amqpConn                                       *amqp.Connection
}

// var _ EventConsumer = (*consumer)(nil)

func NewConsumer(amqpConn *amqp.Connection) EventConsumer {

	return &consumer{
		// exchangeName:   _exchangeName,
		// queueName:      _queueName,
		// bindingKey:     _bindingKey,
		// workerPoolSize: _workerPoolSize,
		amqpConn: amqpConn,
	}
}

func (c *consumer) Configure(opts ...Option) EventConsumer {
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *consumer) StartConsumer(fn worker) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := c.createChannel()
	if err != nil {
		return err
	}
	log.Println("---- Consumer ----")
	msgs, err := ch.Consume(
		c.queueName, // queue
		"Order",     // consumer
		true,        // auto ack
		false,       // exclusive
		false,       // no local
		false,       // no wait
		nil,         // args
	)
	if err != nil {
		log.Println("err consume: ", err.Error())
		return err
	}

	var forever chan struct{}

	for i := 0; i < c.workerPoolSize; i++ {
		go fn(ctx, msgs)
	}
	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever

	return chanErr
}

func (c *consumer) createChannel() (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		c.queueName, // name
		false,       // durable
		false,       // delete when unused
		true,        // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return nil, err
	}

	for _, exchangeName := range c.exchangeName {
		if err = ch.ExchangeDeclare(
			exchangeName, // name
			c.topicType,  // type
			true,         // durable
			false,        // auto-deleted
			false,        // internal
			false,        // no-wait
			nil,          // arguments
		); err != nil {
			return nil, err
		}

		for _, routingKey := range c.RoutingKeys {
			if err = ch.QueueBind(
				q.Name,       // queue name
				routingKey,   // routing key
				exchangeName, // exchange
				_queueNoWait,
				nil,
			); err != nil {
				return nil, err
			}
		}
	}

	return ch, nil
}
