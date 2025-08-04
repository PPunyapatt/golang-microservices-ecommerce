package consumer

type Option func(*Consumer)

func ExchangeName(exchangeName string) Option {
	return func(c *Consumer) {
		c.exchangeName = exchangeName
	}
}

func QueueName(queueName string) Option {
	return func(c *Consumer) {
		c.queueName = queueName
	}
}

func WorkerPoolSize(workerPoolSize int) Option {
	return func(p *Consumer) {
		p.workerPoolSize = workerPoolSize
	}
}
