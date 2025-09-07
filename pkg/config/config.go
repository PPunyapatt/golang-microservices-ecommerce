package config

import (
	"fmt"
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
	postgres_username := os.Getenv("POSTGRES_USERNAME")
	postgres_password := os.Getenv("POSTGRES_PASSWORD")
	postgres_host := os.Getenv("POSTGRES_HOST")
	postgres_port := os.Getenv("POSTGRES_PORT")
	postgres_url := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable", postgres_username, postgres_password, postgres_host, postgres_port)

	// dsn := os.Getenv("POSTGRES_URL")
	rabbitMQ := os.Getenv("RABBITMQ")
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	mongo := os.Getenv("MONGO_URL")
	log.Println("Mongo: ", mongo)
	cfg := &AppConfig{
		ServerPort:  "1024",
		Dsn:         postgres_url,
		RabbitMQUrl: rabbitMQ,
		StripeKey:   stripeKey,
		MongoURL:    mongo,
	}

	return cfg, nil
}
