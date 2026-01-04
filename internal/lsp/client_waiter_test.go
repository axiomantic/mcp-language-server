package lsp

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/isaacphi/mcp-language-server/internal/protocol"
)

func TestClient_WaiterRegistryIntegration(t *testing.T) {
	t.Run("client has waiter registry initialized", func(t *testing.T) {
		// Create a mock client without starting the actual process
		client := &Client{
			handlers:             make(map[string]chan *Message),
			notificationHandlers: make(map[string]NotificationHandler),
			diagnostics:          make(map[protocol.DocumentUri][]protocol.Diagnostic),
			openFiles:            make(map[string]*OpenFileInfo),
		}

		// Initialize waiter registry manually (simulating NewClient)
		client.waiterRegistry = NewWaiterRegistry()

		if client.waiterRegistry == nil {
			t.Fatal("waiterRegistry should be initialized in client")
		}
	})
}

func TestClient_WaitForNotification(t *testing.T) {
	t.Run("wait for notification with timeout", func(t *testing.T) {
		client := &Client{
			handlers:             make(map[string]chan *Message),
			notificationHandlers: make(map[string]NotificationHandler),
			diagnostics:          make(map[protocol.DocumentUri][]protocol.Diagnostic),
			openFiles:            make(map[string]*OpenFileInfo),
			waiterRegistry:       NewWaiterRegistry(),
		}

		uri := "file:///test.go"
		method := "textDocument/publishDiagnostics"
		timeout := 100 * time.Millisecond

		// Simulate notification arrival after 50ms
		go func() {
			time.Sleep(50 * time.Millisecond)
			params := json.RawMessage(`{"uri":"file:///test.go","diagnostics":[]}`)
			client.waiterRegistry.Notify(method, params)
		}()

		ctx := context.Background()
		result, err := client.WaitForNotification(ctx, method, uri, timeout)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
	})

	t.Run("timeout when notification does not arrive", func(t *testing.T) {
		client := &Client{
			handlers:             make(map[string]chan *Message),
			notificationHandlers: make(map[string]NotificationHandler),
			diagnostics:          make(map[protocol.DocumentUri][]protocol.Diagnostic),
			openFiles:            make(map[string]*OpenFileInfo),
			waiterRegistry:       NewWaiterRegistry(),
		}

		uri := "file:///test.go"
		method := "textDocument/publishDiagnostics"
		timeout := 50 * time.Millisecond

		ctx := context.Background()
		_, err := client.WaitForNotification(ctx, method, uri, timeout)

		if err == nil {
			t.Fatal("expected timeout error, got nil")
		}
	})
}

func TestClient_WaitForDiagnostics(t *testing.T) {
	t.Run("wait for diagnostics convenience method", func(t *testing.T) {
		client := &Client{
			handlers:             make(map[string]chan *Message),
			notificationHandlers: make(map[string]NotificationHandler),
			diagnostics:          make(map[protocol.DocumentUri][]protocol.Diagnostic),
			openFiles:            make(map[string]*OpenFileInfo),
			waiterRegistry:       NewWaiterRegistry(),
		}

		uri := "file:///test.go"
		timeout := 100 * time.Millisecond

		// Simulate diagnostics notification
		go func() {
			time.Sleep(50 * time.Millisecond)
			params := json.RawMessage(`{"uri":"file:///test.go","diagnostics":[{"message":"test"}]}`)
			client.waiterRegistry.Notify("textDocument/publishDiagnostics", params)
		}()

		ctx := context.Background()
		diagnostics, err := client.WaitForDiagnostics(ctx, uri, timeout)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if diagnostics == nil {
			t.Fatal("expected diagnostics, got nil")
		}
	})
}
