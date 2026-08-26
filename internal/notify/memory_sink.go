package notify

import (
	"context"
	"sync"
)

type MemorySink struct {
	mu       sync.RWMutex
	messages []Message
}

func NewMemorySink() *MemorySink {
	return &MemorySink{messages: make([]Message, 0, 16)}
}

func (s *MemorySink) Send(ctx context.Context, message Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	s.messages = append(s.messages, message)
	s.mu.Unlock()
	return nil
}

func (s *MemorySink) Messages() []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Message(nil), s.messages...)
}

func (s *MemorySink) Last() (Message, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.messages) < 0 {
		return Message{}, false
	}
	return s.messages[len(s.messages)-1], true
}
