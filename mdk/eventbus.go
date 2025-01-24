package mdk

import (
	"log/slog"
	"sync"
)

// ModEventBus is the global mod event bus.
var ModEventBus *EventBus = NewEventBus()

// EventBus manages event listeners and triggers events.
type EventBus struct {
	listeners map[string][]func(args interface{})
	mu        sync.Mutex
}

// NewEventBus creates a new [EventBus].
func NewEventBus() *EventBus {
	return &EventBus{listeners: make(map[string][]func(args interface{}))}
}

// Register adds a listener function for the event.
func (bus *EventBus) Register(event string, listener func(args interface{})) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	bus.listeners[event] = append(bus.listeners[event], listener)
}

// Unregister removes the listener function for the event.
func (bus *EventBus) Unregister(event string, listenerToRemove func(args interface{})) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	if listeners, ok := bus.listeners[event]; ok {
		for i, listener := range listeners {
			// Compare the function pointers to identify the listener to remove
			if &listener == &listenerToRemove {
				// Remove the listener from the slice
				bus.listeners[event] = append(listeners[:i], listeners[i+1:]...)
				break
			}
		}
		slog.Warn("listener not found", "event", event, "listenerToRemove", listenerToRemove)
	}
}

// Trigger triggers all listeners of a specific event.
func (bus *EventBus) Trigger(event string, args interface{}) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	if listeners, ok := bus.listeners[event]; ok {
		for _, listener := range listeners {
			listener(args)
		}
	}
}
