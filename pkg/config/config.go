package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type AppConfig struct {
	Dsn         string
	RabbitMQUrl string
	StripeKey   string
	MongoURL    string
}

type secret struct {
	Username string
	Password string
	Host     string
	Port     string
}

func SetUpEnv(args ...string) (*AppConfig, error) {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file: ", err.Error())
	// }

	cfg := &AppConfig{}
	for _, key := range args {
		data, err := getSecret(secretPath(key), key)
		if err != nil {
			return nil, fmt.Errorf("error getting %s secret: %w", key, err)
		}
		if data == nil {
			return nil, fmt.Errorf("secret for %s not found", key)
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
			return nil, errors.New("invalid key")
		}
	}

	stripe_data, err := os.ReadFile("/vault/secrets/stripe_secret_key")
	if err != nil {
		log.Println("Error reading STRIPE_SECRET_KEY file: ", err.Error())
	} else {
		parts := strings.SplitN(string(stripe_data), "=", 2)
		stripe_key := parts[1]
		cfg.StripeKey = stripe_key
	}

	return cfg, nil
}

func getSecret(path, key string) (*secret, error) {
	hostPortMap := map[string][2]string{
		"postgres": {"POSTGRES_HOST", "POSTGRES_PORT"},
		"mongodb":  {"MONGO_HOST", "MONGO_PORT"},
		"rabbitmq": {"RABBITMQ_HOST", "RABBITMQ_PORT"},
	}

	hp, ok := hostPortMap[key]
	if !ok {
		return nil, errors.New("invalid key")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var username, password string
	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			switch parts[0] {
			case "username":
				username = parts[1]
			case "password":
				password = parts[1]
			}
		}
	}
	return &secret{
		Username: username,
		Password: password,
		Host:     os.Getenv(hp[0]),
		Port:     os.Getenv(hp[1]),
	}, nil
}

func secretPath(key string) string {
	switch key {
	case "postgres":
		return "/vault/secrets/dbuser"
	case "mongodb":
		return "/vault/secrets/mongouser"
	case "rabbitmq":
		return "/vault/secrets/rabbitmq"
	default:
		return ""
	}
}
