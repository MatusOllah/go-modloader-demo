package mdk

import (
	"errors"
	"maps"
	"sync"
)

var ThingRegistry *Registry[Thing] = NewRegistry[Thing]()

var ErrKeyNotExist error = errors.New("key does not exist")

type RegistryKeyError struct {
	Op  string
	Key string
	Err error
}

func (e *RegistryKeyError) Error() string {
	return e.Op + " " + e.Key + ": " + e.Err.Error()
}

func (e *RegistryKeyError) Unwrap() error {
	return e.Err
}

// Registry is a generic registry for storing objects of any type.
type Registry[T any] struct {
	items map[string]T
	mu    sync.RWMutex
}

// NewRegistry creates a new, empty [Registry].
func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{items: make(map[string]T)}
}

// Register adds an item to the registry with a unique ID.
func (r *Registry[T]) Register(id string, item T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[id] = item
}

// Unregister removes an item by ID.
func (r *Registry[T]) Unregister(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
}

// Get retrieves an item from the registry by ID. If there is an error, it will be of type *RegistryKeyError.
func (r *Registry[T]) Get(id string) (T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[id]
	if !ok {
		return item, &RegistryKeyError{Op: "get", Key: id, Err: ErrKeyNotExist}
	}
	return item, nil
}

// All returns all items in the registry.
func (r *Registry[T]) All() map[string]T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return maps.Clone(r.items)
}
