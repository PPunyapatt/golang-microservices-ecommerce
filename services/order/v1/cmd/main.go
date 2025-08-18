package main

import (
	"config-service"
	"context"
	"net"
	"order/v1/internal/app"
	"order/v1/internal/repository"
	"order/v1/internal/service"
	database "package/Database"
	"package/rabbitmq"
	"package/rabbitmq/constant"
	"package/rabbitmq/consumer"
	"package/rabbitmq/publisher"
	"package/tracer"

	"order/v1/proto/order"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	cfg, err := config.SetUpEnv()
	if err != nil {
		panic(err)
	}

	// ✅ Init tracer
	shutdown := tracer.InitTracer("order-service")
	defer func() { _ = shutdown(context.Background()) }()

	// database connection
	db, err := database.InitDatabase(cfg)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.Gorm.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()
	defer db.Sqlx.Close()

	// RabbitMQ Connection
	conn, err := rabbitmq.NewRabbitMQConnection(cfg.RabbitMQUrl)
	if err != nil {
		panic(err)
	}

	orderRepo := repository.NewOrderRepository(db.Gorm, db.Sqlx)

	orderPublisher, err := publisher.NewPublisher(conn)
	if err != nil {
		panic(err)
	}
	orderPublisher.Configure(
		publisher.TopicType("topic"),
	)

	orderService, orderServiceRPC := service.NewOrderServer(orderRepo, orderPublisher, otel.Tracer("inventory-service"))

	orderQueues := []*constant.Queue{
		{
			Exchange: "inventory.exchange",
			Routing:  "inventory.*",
		},
		{
			Exchange: "payment.exchange",
			Routing:  "payment.*",
		},
	}

	orderConsumer := consumer.NewConsumer(conn)
	orderConsumer.Configure(
		consumer.ExchangeName([]string{
			"inventory.exchange",
			"payment.exchange",
		}),
		consumer.QueueName("order.queue"),
		consumer.RoutingKeys([]string{"payment.*", "inventory.*"}),
		consumer.QueueProperties(orderQueues),
		consumer.WorkerPoolSize(1),
		consumer.TopicType("topic"),
	)

	orderDLQueues := []*constant.Queue{
		{
			Exchange: "inventory.dlx",
			Routing:  "inventory.failed",
		},
		{
			Exchange: "payment.dlx",
			Routing:  "payment.failed",
		},
		{
			Exchange: "order.dlx",
			Routing:  "order.timeout",
		},
	}
	orderDLconsumer := consumer.NewConsumer(conn)
	orderDLconsumer.Configure(
		consumer.ExchangeName([]string{"inventory.dlx", "payment.dlx", "order.dlx"}),
		consumer.RoutingKeys([]string{"inventory.failed", "payment.failed", "order.timeout"}),
		consumer.QueueName("order.dlq"),
		consumer.QueueProperties(orderDLQueues),
		consumer.WorkerPoolSize(1),
		consumer.TopicType("topic"),
	)

	app_ := app.NewWorker(orderService)
	go orderConsumer.StartConsumer(app_.Worker)

	appDlx := app.NewWorkerDeadLetter(orderService)
	go orderDLconsumer.StartConsumer(appDlx.Worker)

	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	// ✅ Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	listener, err := net.Listen("tcp", ":1028")
	if err != nil {
		panic(err)
	}

	order.RegisterOrderServiceServer(s, orderServiceRPC)

	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}
}
