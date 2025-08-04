package consumer

import (
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
	_workerPoolSize = 24
)

type EventConsumer interface {
	Configure(...Option) EventConsumer
	StartConsumer() error
}

type Consumer struct {
	exchangeName, queueName, bindingKey string
	workerPoolSize                      int
	amqpConn                            *amqp.Connection
}

// var _ EventConsumer = (*consumer)(nil)

func NewConsumer(amqpConn *amqp.Connection) *Consumer {

	return &Consumer{
		exchangeName:   _exchangeName,
		queueName:      _queueName,
		bindingKey:     _bindingKey,
		workerPoolSize: _workerPoolSize,
		amqpConn:       amqpConn,
	}
}

func (c *Consumer) Configure(opts ...Option) EventConsumer {
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Consumer) StartConsumer() error {
	return nil
}
