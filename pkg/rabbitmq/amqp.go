package rabbitmq

import "github.com/rabbitmq/amqp091-go"

type AMQPHeaderCarrier amqp091.Table

func (c AMQPHeaderCarrier) Get(key string) string {
	if val, ok := c[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

func (c AMQPHeaderCarrier) Set(key string, value string) {
	c[key] = value
}

func (c AMQPHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}
