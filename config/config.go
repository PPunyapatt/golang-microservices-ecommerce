package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort            string
	Dsn                   string
	AppSecret             string
	TwilioAccountSid      string
	TwilioAuthToken       string
	TwilioFromPhoneNumber string
	RabbitMQUrl           string
}

func SetUpEnv() (*AppConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err.Error())
	}

	dsn := os.Getenv("POSTGRES_URL")
	rabbitMQ := os.Getenv("RABBITMQ")
	log.Println("Dsn: ", dsn)
	cfg := &AppConfig{
		ServerPort:  "1024",
		Dsn:         dsn,
		RabbitMQUrl: rabbitMQ,
	}

	return cfg, nil
}
