package publisher

type Option func(*Publisher)

func ExchangeName(exchangeName []string) Option {
	return func(p *Publisher) {
		p.exchangeName = exchangeName
	}
}

func RoutingKeys(routingKeys []string) Option {
	return func(p *Publisher) {
		p.routingKeys = routingKeys
	}
}

func TopicType(topicType string) Option {
	return func(p *Publisher) {
		p.topicType = topicType
	}
}
