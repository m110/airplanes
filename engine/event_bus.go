package engine

import (
	"fmt"

	"github.com/yohamta/donburi"
)

type Event interface{}

type EventHandler func(w donburi.World, event Event)

type EventBus struct {
	handlers map[Event][]EventHandler

	queue []Event

	debug bool
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: map[Event][]EventHandler{},
		debug:    false,
	}
}

// Process should be called between game ticks to ensure systems don't conflict with each other.
func (e *EventBus) Process(w donburi.World) {
	// The outer loop is needed, because events can trigger more events.
	for len(e.queue) > 0 {
		queue := e.queue
		e.queue = nil
		for _, event := range queue {
			for _, h := range e.handlers[eventName(event)] {
				if e.debug {
					fmt.Printf("%T -> %T\n", event, h)
				}

				h(w, event)
			}
		}
	}
}

func (e *EventBus) Publish(event Event) {
	if e.debug {
		fmt.Printf("Publishing %T\n", event)
	}
	e.queue = append(e.queue, event)
}

func (e *EventBus) Subscribe(event Event, handler EventHandler) {
	if e.debug {
		fmt.Printf("Subscribing %T -> %T\n", event, handler)
	}
	name := eventName(event)
	e.handlers[name] = append(e.handlers[name], handler)
}

func eventName(event Event) string {
	return fmt.Sprintf("%T", event)
}
