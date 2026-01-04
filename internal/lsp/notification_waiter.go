package lsp

import (
	"encoding/json"
	"sync"
)

// NotificationWaiter represents a waiter for a specific notification
type NotificationWaiter struct {
	URI      string
	Ready    chan struct{}
	Result   json.RawMessage
	Received bool
}

// WaiterRegistry manages notification waiters
type WaiterRegistry struct {
	waiters []*NotificationWaiter
	mu      sync.RWMutex
}

// NewWaiterRegistry creates a new waiter registry
func NewWaiterRegistry() *WaiterRegistry {
	return &WaiterRegistry{
		waiters: make([]*NotificationWaiter, 0),
	}
}

// Register adds a waiter to the registry
func (r *WaiterRegistry) Register(waiter *NotificationWaiter) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.waiters = append(r.waiters, waiter)
}

// Unregister removes a waiter from the registry
func (r *WaiterRegistry) Unregister(waiter *NotificationWaiter) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, w := range r.waiters {
		if w == waiter {
			r.waiters = append(r.waiters[:i], r.waiters[i+1:]...)
			break
		}
	}
}

// Notify notifies all matching waiters about a notification
func (r *WaiterRegistry) Notify(method string, params json.RawMessage) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Parse the URI from params for matching
	var paramsMap map[string]interface{}
	if err := json.Unmarshal(params, &paramsMap); err != nil {
		return
	}

	uri, ok := paramsMap["uri"].(string)
	if !ok {
		return
	}

	// Notify matching waiters
	for _, waiter := range r.waiters {
		if waiter.URI == uri && !waiter.Received {
			waiter.Result = params
			waiter.Received = true
			close(waiter.Ready)
		}
	}
}
