package lsp

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/isaacphi/mcp-language-server/internal/protocol"
)

func TestTransport_WaiterNotification(t *testing.T) {
	t.Run("notification triggers waiter registry", func(t *testing.T) {
		// Create a client with waiter registry
		client := &Client{
			handlers:             make(map[string]chan *Message),
			notificationHandlers: make(map[string]NotificationHandler),
			diagnostics:          make(map[protocol.DocumentUri][]protocol.Diagnostic),
			openFiles:            make(map[string]*OpenFileInfo),
			waiterRegistry:       NewWaiterRegistry(),
		}

		// Register a waiter for a specific URI
		waiter := &NotificationWaiter{
			URI:      "file:///test.go",
			Ready:    make(chan struct{}),
			Result:   nil,
			Received: false,
		}
		client.waiterRegistry.Register(waiter)

		// Simulate a notification message being handled
		// This simulates what handleMessages would do when it receives a notification
		msg := &Message{
			JSONRPC: "2.0",
			Method:  "textDocument/publishDiagnostics",
			Params:  json.RawMessage(`{"uri":"file:///test.go","diagnostics":[]}`),
		}

		// Call waiterRegistry.Notify as handleMessages should
		client.waiterRegistry.Notify(msg.Method, msg.Params)

		// Verify waiter was notified
		select {
		case <-waiter.Ready:
			if !waiter.Received {
				t.Error("expected waiter to be marked as received")
			}
			if waiter.Result == nil {
				t.Error("expected waiter result to be set")
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("timeout: waiter was not notified by handleMessages integration")
		}
	})

	t.Run("multiple notifications to different URIs", func(t *testing.T) {
		client := &Client{
			handlers:             make(map[string]chan *Message),
			notificationHandlers: make(map[string]NotificationHandler),
			diagnostics:          make(map[protocol.DocumentUri][]protocol.Diagnostic),
			openFiles:            make(map[string]*OpenFileInfo),
			waiterRegistry:       NewWaiterRegistry(),
		}

		// Register waiters for different URIs
		waiter1 := &NotificationWaiter{
			URI:      "file:///file1.go",
			Ready:    make(chan struct{}),
			Result:   nil,
			Received: false,
		}
		waiter2 := &NotificationWaiter{
			URI:      "file:///file2.go",
			Ready:    make(chan struct{}),
			Result:   nil,
			Received: false,
		}

		client.waiterRegistry.Register(waiter1)
		client.waiterRegistry.Register(waiter2)

		// Notify for file1
		params1 := json.RawMessage(`{"uri":"file:///file1.go","diagnostics":[]}`)
		client.waiterRegistry.Notify("textDocument/publishDiagnostics", params1)

		// Only waiter1 should be notified
		select {
		case <-waiter1.Ready:
			if !waiter1.Received {
				t.Error("waiter1 should be marked as received")
			}
		case <-time.After(50 * time.Millisecond):
			t.Error("waiter1 should have been notified")
		}

		// waiter2 should not be notified yet
		select {
		case <-waiter2.Ready:
			t.Error("waiter2 should not be notified for file1")
		case <-time.After(50 * time.Millisecond):
			// Expected
		}

		// Now notify for file2
		params2 := json.RawMessage(`{"uri":"file:///file2.go","diagnostics":[]}`)
		client.waiterRegistry.Notify("textDocument/publishDiagnostics", params2)

		select {
		case <-waiter2.Ready:
			if !waiter2.Received {
				t.Error("waiter2 should be marked as received")
			}
		case <-time.After(50 * time.Millisecond):
			t.Error("waiter2 should have been notified")
		}
	})
}
