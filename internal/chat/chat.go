package chat

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Chat struct {
	mu        sync.Mutex
	users     map[string]*User
	sessions  map[string]*Session
	waitQueue chan *User
}

func NewChat() *Chat {
	return &Chat{
		users:     make(map[string]*User),
		sessions:  make(map[string]*Session),
		waitQueue: make(chan *User),
	}
}

func (c *Chat) HandleConnection(conn *websocket.Conn) {
	user := NewUser(conn)

	c.mu.Lock()
	c.users[user.ID] = user
	c.mu.Unlock()

	go c.matchmaking(user)
	go user.ReadPump(c)
	go user.WritePump(c)
}

func (c *Chat) matchmaking(user *User) {
	select {
	case partner := <-c.waitQueue:
		session := NewSession(user, partner)

		c.mu.Lock()
		c.sessions[session.ID] = session
		c.mu.Unlock()

		user.SessionID = session.ID
		partner.SessionID = session.ID

		user.SendMsg("Собеседник найден!")
		partner.SendMsg("Собеседник найден!")
	default:
		user.SendMsg("Ожидание собеседника...")
		c.waitQueue <- user
	}
}

func (c *Chat) RemoveSession(sessionID string) {
	fmt.Println("УДАЛИТЬ SESSION", sessionID)
	c.mu.Lock()
	defer c.mu.Unlock()

	if session, ok := c.sessions[sessionID]; ok {
		session.Close()
		delete(c.sessions, sessionID)
	}
}

func (c *Chat) RemoveUser(id string) {
	fmt.Println("УДАЛИТЬ ID", id)
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.users, id)
}
