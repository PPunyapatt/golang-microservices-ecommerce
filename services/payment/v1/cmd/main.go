package main

import (
	"config-service"
	"context"
	"log"
	"package/rabbitmq"
	"package/rabbitmq/publisher"
	"payment/v1/internal/service"
)

func main() {
	cfg, err := config.SetUpEnv()
	if err != nil {
		panic(err)
	}

	// database connection
	// db, err := database.InitDatabase(cfg)
	// if err != nil {
	// 	panic(err)
	// }

	// sqlDB, err := db.Gorm.DB()
	// if err != nil {
	// 	panic(err)
	// }
	// defer sqlDB.Close()
	// defer db.Sqlx.Close()

	// RabbitMQ Connection
	conn, err := rabbitmq.NewRabbitMQConnection(cfg.RabbitMQUrl)
	if err != nil {
		panic(err)
	}

	publisher := publisher.NewPublisher(conn)
	paymentService := service.NewPaymentService(cfg.StripeKey, publisher)

	if err = paymentService.ProcessPayment(context.Background(), 1, 543.21, "9e49ca9b-a4e9-4528-af11-1978b23c185f"); err != nil {
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

}
