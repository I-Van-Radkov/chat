package server

import (
	"net/http"

	"github.com/I-Van-Radkov/chat/internal/chat"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	router *gin.Engine
	chat   *chat.Chat
}

func NewServer(chat *chat.Chat) *Server {
	router := gin.Default()

	router.Use(corsMiddleware())

	server := &Server{
		router: router,
		chat:   chat,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.router.GET("/ws", s.handlerChat)
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		c.Next()
	}
}
