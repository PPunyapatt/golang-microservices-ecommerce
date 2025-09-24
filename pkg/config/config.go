package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Dsn                 string
	RabbitMQUrl         string
	StripeKey           string
	StripeWebhookSecret string
	MongoURL            string
	JwtSecret           string
}

type secret struct {
	Username            string
	Password            string
	Host                string
	Port                string
	StripeSecret        string
	StripeWebhookSecret string
}

func SetUpEnv(args ...string) (*AppConfig, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, skipping...")
	}
	env := os.Getenv("ENVIRONMENT")

	cfg := &AppConfig{}
	for _, key := range args {
		switch env {
		case "prod":
			getVault(cfg, key)
		case "dev":
			getEnv(cfg, key)
		}
	}

	return cfg, nil
}

func getVault(cfg *AppConfig, key string) error {
	data, err := getSecret(secretPath(key), key)
	if err != nil {
		return fmt.Errorf("error getting %s secret: %w", key, err)
	}
	if data == nil {
		return fmt.Errorf("secret for %s not found", key)
	}

	switch key {
	case "postgres":
		cfg.Dsn = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/postgres?sslmode=disable",
			data.Username, data.Password, data.Host, data.Port,
		)
	case "mongodb":
		cfg.MongoURL = fmt.Sprintf(
			"mongodb://%s:%s@%s:%s/ecommerce?authSource=admin",
			data.Username, data.Password, data.Host, data.Port,
		)
	case "rabbitmq":
		cfg.RabbitMQUrl = fmt.Sprintf(
			"amqp://%s:%s@%s:%s/",
			data.Username, data.Password, data.Host, data.Port,
		)
	default:
		return errors.New("invalid key")
	}

	stripe_data, err := os.ReadFile("/vault/secrets/stripe_secret_key")
	if err != nil {
		log.Println("Error reading STRIPE_SECRET_KEY file: ", err.Error())
	} else {
		parts := strings.SplitN(string(stripe_data), "=", 2)
		cfg.StripeKey = parts[1]
	}

	jwtSecret_data, err := os.ReadFile("/vault/secrets/jwt")
	if err != nil {
		log.Fatal("Error reading jwt secret file: ", err.Error())
	} else {
		parts := strings.SplitN(string(jwtSecret_data), "=", 2)
		cfg.JwtSecret = parts[1]
	}

	return nil
}

func getEnv(cfg *AppConfig, key string) error {
	switch key {
	case "postgres":
		cfg.Dsn = os.Getenv("POSTGRES_URL")
	case "mongodb":
		cfg.MongoURL = os.Getenv("MONGO_URL")
	case "rabbitmq":
		cfg.RabbitMQUrl = os.Getenv("RABBITMQ")
	default:
		return errors.New("invalid key")
	}

	cfg.StripeKey = os.Getenv("STRIPE_SECRET_KEY")
	cfg.StripeWebhookSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")
	cfg.JwtSecret = os.Getenv("JWT_SECRET")

	return nil
}

func getSecret(path, key string) (*secret, error) {
	hostPortMap := map[string][2]string{
		"postgres": {"POSTGRES_HOST", "POSTGRES_PORT"},
		"mongodb":  {"MONGO_HOST", "MONGO_PORT"},
		"rabbitmq": {"RABBITMQ_HOST", "RABBITMQ_PORT"},
	}

	secret := &secret{}
	hp, ok := hostPortMap[key]
	if ok {
		secret.Host = os.Getenv(hp[0])
		secret.Port = os.Getenv(hp[1])
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			switch parts[0] {
			case "username":
				secret.Username = parts[1]
			case "password":
				secret.Password = parts[1]
			case "secret_key":
				secret.StripeSecret = parts[1]
			case "webhook":
				secret.StripeWebhookSecret = parts[1]
			}
		}
	}

	return secret, nil
}

func secretPath(key string) string {
	switch key {
	case "postgres":
		return "/vault/secrets/dbuser"
	case "mongodb":
		return "/vault/secrets/mongouser"
	case "rabbitmq":
		return "/vault/secrets/rabbitmq"
	case "stripe-key":
		return "/vault/secrets/stripe-key"
	default:
		return ""
	}
}

func ReadVaultSecret(key string) (*secret, error) {
	data := &secret{}
	var err error

	env := os.Getenv("ENVIRONMENT")
	switch env {
	case "prod":
		data, err = getSecret(secretPath(key), key)
		if err != nil {
			return nil, fmt.Errorf("error getting %s secret: %w", key, err)
		}
	case "dev":
		data.StripeSecret = os.Getenv("STRIPE_SECRET_KEY")
		data.StripeWebhookSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")
	}
	return data, nil
}
