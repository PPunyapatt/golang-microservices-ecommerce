package config

import (
	"fmt"
	"log"
	"os"
	"strings"
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

	var postgres_username, postgres_password string
	postgres_data, err := os.ReadFile("/vault/secrets/dbuser")
	if err != nil {
		// log.Fatal("Error reading dbuser file: ", err.Error())
		log.Println("Error reading dbuser file: ", err.Error())
	}
	log.Println("postgres_data: ", string(postgres_data))
	for _, line := range strings.Split(string(postgres_data), "\n") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			switch parts[0] {
			case "username":
				postgres_username = parts[1]
			case "password":
				postgres_password = parts[1]
			}
		}
	}

	postgres_host := os.Getenv("POSTGRES_HOST")
	postgres_port := os.Getenv("POSTGRES_PORT")
	postgres_url := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable", postgres_username, postgres_password, postgres_host, postgres_port)

	// dsn := os.Getenv("POSTGRES_URL")
	rabbitMQ := os.Getenv("RABBITMQ")
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	mongo := os.Getenv("MONGO_URL")
	cfg := &AppConfig{
		ServerPort:  "1024",
		Dsn:         postgres_url,
		RabbitMQUrl: rabbitMQ,
		StripeKey:   stripeKey,
		MongoURL:    mongo,
	}

	return cfg, nil
}
