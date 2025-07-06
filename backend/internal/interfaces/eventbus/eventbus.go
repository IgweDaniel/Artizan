package interfaces

type Event struct {
	EventType string
	Data      interface{}
}

type HandlerFunc func(Event)

type EventBus interface {
	Subscribe(event string, handler HandlerFunc)
	Publish(event string, data interface{})
}
