package config

import (
	"context"
	"sync"
	"time"
)

const DefaultChannelBuffer = 32

type Event struct {
	Type      string    `json:"type"`
	Payload   any       `json:"payload,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type EventHub struct {
	events      chan Event
	register    chan chan Event
	unregister  chan chan Event
	subscribers map[chan Event]struct{}
	once        sync.Once
}

func NewEventHub(buffer int) *EventHub {
	if buffer <= 0 {
		buffer = DefaultChannelBuffer
	}

	return &EventHub{
		events:      make(chan Event, buffer),
		register:    make(chan chan Event, buffer),
		unregister:  make(chan chan Event, buffer),
		subscribers: make(map[chan Event]struct{}),
	}
}

func (hub *EventHub) Publish(eventType string, payload any) bool {
	event := Event{
		Type:      eventType,
		Payload:   payload,
		Timestamp: time.Now().UTC(),
	}

	select {
	case hub.events <- event:
		return true
	default:
		return false
	}
}

func (hub *EventHub) Subscribe(ctx context.Context) <-chan Event {
	subscriber := make(chan Event, DefaultChannelBuffer)

	select {
	case hub.register <- subscriber:
	case <-ctx.Done():
		close(subscriber)
		return subscriber
	}

	go func() {
		<-ctx.Done()

		select {
		case hub.unregister <- subscriber:
		default:
		}
	}()

	return subscriber
}

func (hub *EventHub) Run(ctx context.Context) {
	hub.once.Do(func() {
		for {
			select {
			case <-ctx.Done():
				for subscriber := range hub.subscribers {
					close(subscriber)
					delete(hub.subscribers, subscriber)
				}
				return
			case subscriber := <-hub.register:
				hub.subscribers[subscriber] = struct{}{}
			case subscriber := <-hub.unregister:
				if _, ok := hub.subscribers[subscriber]; ok {
					close(subscriber)
					delete(hub.subscribers, subscriber)
				}
			case event := <-hub.events:
				for subscriber := range hub.subscribers {
					select {
					case subscriber <- event:
					default:
					}
				}
			}
		}
	})
}
