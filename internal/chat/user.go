package chat

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type User struct {
	ID        string
	Conn      *websocket.Conn
	SessionID string
	Send      chan []byte
	closeChan chan struct{}
	once      sync.Once
}

func NewUser(conn *websocket.Conn) *User {
	return &User{
		ID:        generateUUID(),
		Conn:      conn,
		Send:      make(chan []byte, 256),
		closeChan: make(chan struct{}),
	}
}

func generateUUID() string {
	return uuid.NewString()
}

func (u *User) ReadPump(c *Chat) {
	defer func() {
		c.RemoveUser(u.ID)
		c.RemoveSession(u.SessionID)
		u.Stop()
	}()

	for {
		_, msg, err := u.Conn.ReadMessage()
		if err != nil {
			break
		}

		c.mu.Lock()
		session, exists := c.sessions[u.SessionID]
		c.mu.Unlock()

		if string(msg) == "!с" {
			break
		}

		if exists {
			session.Broadcast(u.ID, msg)
		}
	}
}

func (u *User) WritePump(c *Chat) {
	for {
		select {
		case msg, ok := <-u.Send:
			if !ok {
				return
			}

			if err := u.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-u.closeChan:
			return
		}
	}
}

func (u *User) Stop() {
	u.once.Do(func() {
		fmt.Println("ОСТАНОВИТЬ", u.ID)
		close(u.closeChan)
		u.Conn.Close()
	})
}

func (u *User) SendMsg(msg string) {
	u.Send <- []byte(msg)
}
