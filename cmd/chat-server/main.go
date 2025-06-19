package main

import (
	"log"
	"net/http"

	"github.com/I-Van-Radkov/chat.git/internal/chat"
	"github.com/I-Van-Radkov/chat.git/internal/server"
)

const (
	addr = ":8080"
)

func main() {
	chatService := chat.NewChat()

	srv := server.NewServer(chatService)

	http.HandleFunc("/ws", srv.Handler)

	log.Printf("Сервер запущен на %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Ошибка запуска сервера")
	}
}
