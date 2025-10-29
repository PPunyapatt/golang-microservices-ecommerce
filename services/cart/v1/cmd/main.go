package main

import (
	"cart/v1/internal/server"
	"log/slog"
	"os"
	"package/config"
)

func main() {
	cfg, err := config.SetUpEnv("mongodb", "rabbitmq")
	if err != nil {
		panic(err)
	}

	s := server.NewServer(cfg)
	if err := s.Run(); err != nil {
		slog.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}
