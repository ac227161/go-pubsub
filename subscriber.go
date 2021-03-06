package pubsub

import (
	"sync"
)

// Subscribers ...
type Subscribers map[string]*Subscriber

// Subscriber ...
type Subscriber struct {
	id        string
	messages  chan *Message
	createdAt int64
	destroyed bool
	lock      *sync.RWMutex
	topics    map[string]bool
}

// GetID return the subscriber id
func (s *Subscriber) GetID() string {
	return s.id
}

// GetCreatedAt return `time.Time` of the creation time
func (s *Subscriber) GetCreatedAt() int64 {
	return s.createdAt
}

// GetTopics return slice of subscriber topics
func (s *Subscriber) GetTopics() []string {
	topics := []string{}
	for topic := range s.topics {
		topics = append(topics, topic)
	}
	return topics
}

// GetMessages returns a channel of *Message to listen on
func (s *Subscriber) GetMessages() <-chan *Message {
	return s.messages
}

// Signal sends a message to subscriber
func (s *Subscriber) Signal(m *Message) *Subscriber {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if !s.destroyed {
		s.messages <- m
	}
	return s
}

// close the underlying channels/resources
func (s *Subscriber) destroy() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.destroyed = true
	close(s.messages)
}
