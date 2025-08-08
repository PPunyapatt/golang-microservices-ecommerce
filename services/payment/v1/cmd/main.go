package main

import (
	"config-service"
	"context"
	"log"
	"net"
	"package/rabbitmq"
	"package/rabbitmq/publisher"
	"package/tracer"
	"payment/v1/internal/service"
	"payment/v1/proto/payment"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	cfg, err := config.SetUpEnv()
	if err != nil {
		panic(err)
	}

	shutdown := tracer.InitTracer("payment-service")
	defer func() { _ = shutdown(context.Background()) }()

	// RabbitMQ Connection
	conn, err := rabbitmq.NewRabbitMQConnection(cfg.RabbitMQUrl)
	if err != nil {
		panic(err)
	}

	publisher := publisher.NewPublisher(conn)
	paymentService, paymentServiceRPC := service.NewPaymentService(cfg.StripeKey, publisher)

	if err = paymentService.ProcessPayment(context.Background(), 1, 1500.78, "9e49ca9b-a4e9-4528-af11-1978b23c185f"); err != nil {
		log.Println("Err process payment: ", err.Error())
		panic(err)
	}
	// paymentConsumer := consumer.NewConsumer(conn)
	// paymentConsumer.Configure(
	// 	consumer.ExchangeName("payment.exchange"),
	// 	consumer.QueueName("payment"),
	// 	consumer.RoutingKeys([]string{"payment.*"}),
	// 	consumer.WorkerPoolSize(2),
	// 	consumer.TopicType("topic"),
	// )

	// app := app.NewWorker(paymentService)

	// go paymentConsumer.StartConsumer(app.Worker)

	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	listener, err := net.Listen("tcp", ":1029")
	if err != nil {
		panic(err)
	}

	payment.RegisterPaymentServiceServer(s, paymentServiceRPC)

	// ✅ Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}
}
