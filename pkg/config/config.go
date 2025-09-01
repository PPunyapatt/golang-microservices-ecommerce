package config

import (
	"log"
	"os"
)

type AppConfig struct {
	ServerPort            string
	Dsn                   string
	AppSecret             string
	TwilioAccountSid      string
	TwilioAuthToken       string
	TwilioFromPhoneNumber string
	RabbitMQUrl           string
	StripeKey             string
	MongoURL              string
}

func SetUpEnv() (*AppConfig, error) {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file: ", err.Error())
	// }

	dsn := os.Getenv("POSTGRES_URL")
	rabbitMQ := os.Getenv("RABBITMQ")
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	mongo := os.Getenv("MONGO_URL")
	log.Println("Dsn: ", dsn)
	log.Println("Mongo: ", mongo)
	cfg := &AppConfig{
		ServerPort:  "1024",
		Dsn:         dsn,
		RabbitMQUrl: rabbitMQ,
		StripeKey:   stripeKey,
		MongoURL:    mongo,
	}

	return cfg, nil
}
