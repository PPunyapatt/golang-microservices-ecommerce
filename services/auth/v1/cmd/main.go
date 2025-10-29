package main

import (
	"auth-service/v1/internal/server"
	"log/slog"
	_ "net/http/pprof"
	"os"
	"package/config"
)

func main() {
	cfg, err := config.SetUpEnv("postgres")
	if err != nil {
		panic(err)
	}

	s := server.NewServer(cfg)
	if err := s.Run(); err != nil {
		slog.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}
