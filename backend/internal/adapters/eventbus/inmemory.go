package eventbus

import (
	"sync"

	interfaces "github.com/igwedaniel/artizan/internal/interfaces/eventbus"
)

type inMemoryBus struct {
	handlers map[string][]interfaces.HandlerFunc
	mu       sync.RWMutex
}

func New() interfaces.EventBus {
	return &inMemoryBus{
		handlers: make(map[string][]interfaces.HandlerFunc),
	}
}

func (b *inMemoryBus) Subscribe(eventType string, handler interfaces.HandlerFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *inMemoryBus) Publish(eventType string, data interface{}) {
	b.mu.RLock()
	handlers := b.handlers[eventType]
	b.mu.RUnlock()
	event := interfaces.Event{
		EventType: eventType,
		Data:      data,
	}
	for _, h := range handlers {
		go h(event)
	}
}
