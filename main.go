package main

import (
	"fmt"
	"log"

	"github.com/jcoelho93/irc/internal/server"
)

func main() {
	port := ":8080"
	fmt.Println("Starting server on port", port)
	server := server.NewInternetRelayChatServer(port)
	err := server.Start()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	fmt.Println("Server stopped")
}
