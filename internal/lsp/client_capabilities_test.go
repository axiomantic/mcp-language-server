package lsp

import (
	"testing"

	"github.com/isaacphi/mcp-language-server/internal/protocol"
)

func TestClient_GetCapabilities(t *testing.T) {
	t.Run("returns nil when capabilities not set", func(t *testing.T) {
		client := &Client{
			handlers:             make(map[string]chan *Message),
			notificationHandlers: make(map[string]NotificationHandler),
			diagnostics:          make(map[protocol.DocumentUri][]protocol.Diagnostic),
			openFiles:            make(map[string]*OpenFileInfo),
			waiterRegistry:       NewWaiterRegistry(),
		}

		caps := client.GetCapabilities()
		if caps != nil {
			t.Error("expected nil capabilities when not initialized")
		}
	})

	t.Run("returns capabilities after initialization", func(t *testing.T) {
		client := &Client{
			handlers:             make(map[string]chan *Message),
			notificationHandlers: make(map[string]NotificationHandler),
			diagnostics:          make(map[protocol.DocumentUri][]protocol.Diagnostic),
			openFiles:            make(map[string]*OpenFileInfo),
			waiterRegistry:       NewWaiterRegistry(),
		}

		// Simulate setting capabilities after initialization
		testCaps := &protocol.ServerCapabilities{
			TextDocumentSync: &protocol.Or_ServerCapabilities_textDocumentSync{
				Value: &protocol.TextDocumentSyncOptions{
					OpenClose: true,
					Change:    protocol.Full,
				},
			},
		}

		// Set capabilities (simulating what InitializeLSPClient would do)
		client.capabilities = testCaps

		// Get capabilities
		caps := client.GetCapabilities()
		if caps == nil {
			t.Fatal("expected capabilities to be returned")
		}

		// Verify it's the same instance
		if caps != testCaps {
			t.Error("expected same capabilities instance")
		}
	})

	t.Run("capabilities have expected fields", func(t *testing.T) {
		client := &Client{
			handlers:             make(map[string]chan *Message),
			notificationHandlers: make(map[string]NotificationHandler),
			diagnostics:          make(map[protocol.DocumentUri][]protocol.Diagnostic),
			openFiles:            make(map[string]*OpenFileInfo),
			waiterRegistry:       NewWaiterRegistry(),
		}

		testCaps := &protocol.ServerCapabilities{
			HoverProvider: &protocol.Or_ServerCapabilities_hoverProvider{
				Value: true,
			},
			DefinitionProvider: &protocol.Or_ServerCapabilities_definitionProvider{
				Value: true,
			},
			CompletionProvider: &protocol.CompletionOptions{
				TriggerCharacters: []string{".", ":"},
			},
		}

		client.capabilities = testCaps
		caps := client.GetCapabilities()

		if caps.HoverProvider == nil {
			t.Error("expected HoverProvider to be set")
		}
		if caps.DefinitionProvider == nil {
			t.Error("expected DefinitionProvider to be set")
		}
		if caps.CompletionProvider == nil {
			t.Error("expected CompletionProvider to be set")
		}
	})
}
