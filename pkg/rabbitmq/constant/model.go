package constant

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	Exchange string
	Routing  string
}

type ConsumerConfig struct {
	QueueName      string
	Bindings       []*Queue
	DeadLetter     *Queue
	WorkerPoolSize int
	Handler        func(ctx context.Context, messages amqp091.Delivery)
	StartWorker    bool // สำหรับ delay queue ที่ไม่ run worker
}
