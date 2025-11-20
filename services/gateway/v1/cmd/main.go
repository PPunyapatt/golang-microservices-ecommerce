package main

import (
	"gateway/v1/internal/server"
	"log/slog"
	"os"
	"package/config"
)

func main() {
	cfg, err := config.SetUpEnv("postgres", "rabbitmq")
	if err != nil {
		panic(err)
	}

	if err := server.Run(cfg); err != nil {
		slog.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}
