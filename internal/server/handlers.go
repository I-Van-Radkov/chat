package server

import (
	"log"

	"github.com/gin-gonic/gin"
)

func (s *Server) handlerChat(c *gin.Context) {
	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	s.chat.HandleConnection(conn)
}
