package main

import (
	"log/slog"

	"github.com/jcoelho93/irc/internal/server"
)

func main() {
	port := ":8080"
	server := server.NewInternetRelayChatServer(port)
	err := server.Start()
	if err != nil {
		slog.Error("Failed to start IRC server", "error", err)
	}
	slog.Info("IRC server stopped")
}
