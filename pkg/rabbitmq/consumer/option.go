package consumer

import "package/rabbitmq/constant"

type Option func(*consumer)

func ExchangeName(exchangeName []string) Option {
	return func(c *consumer) {
		c.exchangeName = exchangeName
	}
}

func QueueName(queueName string) Option {
	return func(c *consumer) {
		c.queueName = queueName
	}
}

func RoutingKeys(queueNames []string) Option {
	return func(c *consumer) {
		c.RoutingKeys = queueNames
	}
}

func WorkerPoolSize(workerPoolSize int) Option {
	return func(p *consumer) {
		p.workerPoolSize = workerPoolSize
	}
}

func TopicType(topicType string) Option {
	return func(p *consumer) {
		p.topicType = topicType
	}
}

func QueueProperties(queues []*constant.Queue) Option {
	return func(p *consumer) {
		p.queues = queues
	}
}

func QueueDeadLetter(queue *constant.Queue) Option {
	return func(p *consumer) {
		p.queueDeadLetter = queue
	}
}
