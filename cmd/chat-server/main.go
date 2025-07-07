package main

import (
	"log"

	"github.com/I-Van-Radkov/chat/internal/chat"
	"github.com/I-Van-Radkov/chat/internal/server"
)

const (
	addr = ":8080"
)

func main() {
	chatService := chat.NewChat()

	srv := server.NewServer(chatService)

	if err := srv.Start(addr); err != nil {
		log.Fatal("Ошибка запуска сервера")
	}
}
