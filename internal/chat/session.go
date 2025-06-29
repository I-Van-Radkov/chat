package chat

import (
	"sync"
	"time"
)

type Session struct {
	ID        string
	User1     *User
	User2     *User
	Messages  []Message
	CreatedAt time.Time
	mu        sync.Mutex
	once      sync.Once
}

type Message struct {
	SenderID  string
	Content   string
	Timestamp time.Time
}

func NewSession(user1, user2 *User) *Session {
	return &Session{
		ID:        generateUUID(),
		User1:     user1,
		User2:     user2,
		CreatedAt: time.Now(),
	}
}

func (s *Session) Broadcast(sendderId string, msg []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Messages = append(s.Messages, Message{
		SenderID:  sendderId,
		Content:   string(msg),
		Timestamp: time.Now(),
	})

	target := s.User1
	if sendderId == target.ID {
		target = s.User2
	}

	select {
	case target.Send <- msg:
	default:
		s.Close()
	}

}

func (s *Session) Close() {
	s.once.Do(func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		if s.User1 != nil {
			s.User1.Disconnect()
		}

		if s.User2 != nil {
			s.User2.Disconnect()
		}

		s.User1 = nil
		s.User2 = nil
	})
}
