package lsp

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNotificationWaiter(t *testing.T) {
	t.Run("waiter is created with ready channel", func(t *testing.T) {
		waiter := &NotificationWaiter{
			URI:      "file:///test.go",
			Ready:    make(chan struct{}),
			Result:   nil,
			Received: false,
		}

		if waiter.URI != "file:///test.go" {
			t.Errorf("expected URI file:///test.go, got %s", waiter.URI)
		}
		if waiter.Received {
			t.Errorf("expected Received to be false, got true")
		}
		if waiter.Result != nil {
			t.Errorf("expected Result to be nil, got %v", waiter.Result)
		}
	})
}

func TestWaiterRegistry_RegisterUnregister(t *testing.T) {
	registry := NewWaiterRegistry()

	t.Run("register and retrieve waiter", func(t *testing.T) {
		waiter := &NotificationWaiter{
			URI:      "file:///test.go",
			Ready:    make(chan struct{}),
			Result:   nil,
			Received: false,
		}

		registry.Register(waiter)

		// Verify registration by notifying
		testParams := json.RawMessage(`{"uri":"file:///test.go"}`)
		registry.Notify("textDocument/publishDiagnostics", testParams)

		select {
		case <-waiter.Ready:
			if !waiter.Received {
				t.Error("expected Received to be true after notification")
			}
			if waiter.Result == nil {
				t.Error("expected Result to be set after notification")
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("timeout waiting for notification")
		}
	})

	t.Run("unregister waiter", func(t *testing.T) {
		waiter := &NotificationWaiter{
			URI:      "file:///test2.go",
			Ready:    make(chan struct{}),
			Result:   nil,
			Received: false,
		}

		registry.Register(waiter)
		registry.Unregister(waiter)

		// After unregister, notification should not trigger
		testParams := json.RawMessage(`{"uri":"file:///test2.go"}`)
		registry.Notify("textDocument/publishDiagnostics", testParams)

		select {
		case <-waiter.Ready:
			t.Error("waiter should not be notified after unregister")
		case <-time.After(50 * time.Millisecond):
			// Expected - no notification
		}
	})
}

func TestWaiterRegistry_Notify(t *testing.T) {
	registry := NewWaiterRegistry()

	t.Run("notify matching URI", func(t *testing.T) {
		waiter := &NotificationWaiter{
			URI:      "file:///match.go",
			Ready:    make(chan struct{}),
			Result:   nil,
			Received: false,
		}
		registry.Register(waiter)
		defer registry.Unregister(waiter)

		testParams := json.RawMessage(`{"uri":"file:///match.go","diagnostics":[]}`)
		registry.Notify("textDocument/publishDiagnostics", testParams)

		select {
		case <-waiter.Ready:
			if !waiter.Received {
				t.Error("expected Received to be true")
			}
			if string(waiter.Result) != string(testParams) {
				t.Errorf("expected Result to match params, got %s", string(waiter.Result))
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("timeout waiting for notification")
		}
	})

	t.Run("do not notify non-matching URI", func(t *testing.T) {
		waiter := &NotificationWaiter{
			URI:      "file:///different.go",
			Ready:    make(chan struct{}),
			Result:   nil,
			Received: false,
		}
		registry.Register(waiter)
		defer registry.Unregister(waiter)

		testParams := json.RawMessage(`{"uri":"file:///other.go","diagnostics":[]}`)
		registry.Notify("textDocument/publishDiagnostics", testParams)

		select {
		case <-waiter.Ready:
			t.Error("waiter should not be notified for non-matching URI")
		case <-time.After(50 * time.Millisecond):
			// Expected - no notification
		}
	})
}

func TestWaiterRegistry_ConcurrentAccess(t *testing.T) {
	registry := NewWaiterRegistry()

	// Test concurrent registration
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			waiter := &NotificationWaiter{
				URI:      "file:///test.go",
				Ready:    make(chan struct{}),
				Result:   nil,
				Received: false,
			}
			registry.Register(waiter)
			registry.Unregister(waiter)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}
