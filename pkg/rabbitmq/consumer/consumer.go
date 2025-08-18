package consumer

import (
	"context"
	"log"
	"package/rabbitmq/constant"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type worker func(ctx context.Context, messages <-chan amqp.Delivery)
type EventConsumer interface {
	Configure(...Option) EventConsumer
	StartConsumer(worker) error
}

type consumer struct {
	bindingKey, consumerName, topicType, queueName string
	workerPoolSize                                 int
	RoutingKeys, exchangeName                      []string
	amqpConn                                       *amqp.Connection
	deadLetter                                     bool
	queues                                         []*constant.Queue
	queueDeadLetter                                *constant.Queue
}

// var _ EventConsumer = (*consumer)(nil)

func NewConsumer(amqpConn *amqp.Connection, dl ...bool) EventConsumer {

	return &consumer{
		amqpConn:   amqpConn,
		deadLetter: dl[0],
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
	var wg sync.WaitGroup
	for _, queueName := range c.queueName {
		msgs, err := ch.Consume(
			queueName, // queue
			"",        // consumer
			false,     // auto ack
			false,     // exclusive
			false,     // no local
			false,     // no wait
			nil,       // args
		)
		if err != nil {
			log.Println("err consume: ", err.Error())
			return err
		}

		for i := 0; i < c.workerPoolSize; i++ {
			wg.Add(1)
			go func(messages <-chan amqp.Delivery) {
				defer wg.Done()
				fn(ctx, messages)
			}(msgs)
		}

	}
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	// รอจนกว่า connection ปิด
	go func() {
		err := <-ch.NotifyClose(make(chan *amqp.Error))
		if err != nil {
			log.Println("channel closed:", err)
		}
		cancel()
	}()

	<-ctx.Done()
	wg.Wait()
	return nil
}

func (c *consumer) createChannel() (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()
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
	}

	for _, queueName := range c.queueName {
		q, err := ch.QueueDeclare(
			queueName, // name
			false,     // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			return nil, err
		}

		for _, exchangeName := range c.exchangeName {
			for _, routingKey := range c.RoutingKeys {
				if err = ch.QueueBind(
					q.Name,       // queue name
					routingKey,   // routing key
					exchangeName, // exchange
					false,
					nil,
				); err != nil {
					return nil, err
				}
			}
		}
	}

	for _, queue := range c.queues {
		if err = ch.ExchangeDeclare(
			queue.Exchange, // name
			c.topicType,    // type
			true,           // durable
			false,          // auto-deleted
			false,          // internal
			false,          // no-wait
			nil,            // arguments
		); err != nil {
			return nil, err
		}

		var args *amqp.Table
		if c.queueDeadLetter != nil {
			args = &amqp.Table{
				"x-dead-letter-exchange":    c.queueDeadLetter.Exchange,
				"x-dead-letter-routing-key": c.queueDeadLetter.Routing,
			}
		}

		q, err := ch.QueueDeclare(
			c.queueName, // name
			false,       // durable
			false,       // delete when unused
			false,       // exclusive
			false,       // no-wait
			*args,       // arguments
		)
		if err != nil {
			return nil, err
		}

		if err = ch.QueueBind(
			q.Name,         // queue name
			queue.Routing,  // routing key
			queue.Exchange, // exchange
			false,
			nil,
		); err != nil {
			return nil, err
		}
	}

	return ch, nil
}
