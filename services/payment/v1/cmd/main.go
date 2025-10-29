package main

import (
	"log/slog"
	"os"
	"package/config"
	"payment/v1/internal/server"
)

func main() {
	cfg, err := config.SetUpEnv("postgres", "rabbitmq")
	if err != nil {
		panic(err)
	}

	s := server.NewServer(cfg)
	if err := s.Run(); err != nil {
		slog.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}
