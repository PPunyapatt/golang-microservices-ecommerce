package publisher

type Option func(*Publisher)

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
