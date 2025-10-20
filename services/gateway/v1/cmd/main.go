package main

import (
	"gateway/v1/internal/server"
	"log/slog"
	"os"
)

func main() {
	if err := server.Run(); err != nil {
		slog.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}
