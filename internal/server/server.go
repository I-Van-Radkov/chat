package server

import (
	"context"
	"log"
	"net/http"

	"github.com/I-Van-Radkov/chat.git/internal/chat"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	chat *chat.Chat
}

func NewServer(chat *chat.Chat) *Server {
	return &Server{
		chat: chat,
	}
}

func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	ctx := context.WithValue(context.Background(), "userId", userId)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	s.chat.HandleConnection(ctx, conn)
}

/*type Server struct {
	router     *gin.Engine
	chat       *chat.Chat
	httpServer *http.Server
}

func NewServer(chat *chat.Chat) *Server {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		corsMiddleware(),
	)

	server := &Server{
		router: router,
		chat:   chat,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {

}

func (s *Server) handleWebSocket(c *gin.Context) {
}

func (s *Server) Start() error {
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		c.Next()
	}
}
*/
