package integrationtests

import (
	"testing"

	"github.com/isaacphi/mcp-language-server/internal/protocol"
)

// MockServerCapabilities creates test capabilities with specified support
type MockServerCapabilities struct {
	SupportedCapabilities []string
}

// ToProtocol converts mock capabilities to actual ServerCapabilities struct
func (m *MockServerCapabilities) ToProtocol() *protocol.ServerCapabilities {
	caps := &protocol.ServerCapabilities{}

	for _, cap := range m.SupportedCapabilities {
		switch cap {
		case "definition":
			caps.DefinitionProvider = &protocol.Or_ServerCapabilities_definitionProvider{
				Value: true,
			}
			caps.WorkspaceSymbolProvider = &protocol.Or_ServerCapabilities_workspaceSymbolProvider{
				Value: true,
			}
		case "references":
			caps.ReferencesProvider = &protocol.Or_ServerCapabilities_referencesProvider{
				Value: true,
			}
		case "hover":
			caps.HoverProvider = &protocol.Or_ServerCapabilities_hoverProvider{
				Value: true,
			}
		case "rename":
			caps.RenameProvider = true
		case "code_actions":
			caps.CodeActionProvider = true
		case "signature_help":
			caps.SignatureHelpProvider = &protocol.SignatureHelpOptions{}
		case "document_symbols":
			caps.DocumentSymbolProvider = &protocol.Or_ServerCapabilities_documentSymbolProvider{
				Value: true,
			}
		case "call_hierarchy":
			caps.CallHierarchyProvider = &protocol.Or_ServerCapabilities_callHierarchyProvider{
				Value: true,
			}
		case "code_lens":
			caps.CodeLensProvider = &protocol.CodeLensOptions{}
		}
	}

	return caps
}

// TestMockServerCapabilities verifies the mock helper works correctly
func TestMockServerCapabilities(t *testing.T) {
	tests := []struct {
		name          string
		capabilities  []string
		checkFunc     func(*protocol.ServerCapabilities) bool
		checkName     string
		shouldSupport bool
	}{
		{
			name:          "definition capability",
			capabilities:  []string{"definition"},
			checkFunc:     func(c *protocol.ServerCapabilities) bool { return c.DefinitionProvider != nil && c.DefinitionProvider.Value != nil },
			checkName:     "DefinitionProvider",
			shouldSupport: true,
		},
		{
			name:          "no definition capability",
			capabilities:  []string{"hover"},
			checkFunc:     func(c *protocol.ServerCapabilities) bool { return c.DefinitionProvider != nil && c.DefinitionProvider.Value != nil },
			checkName:     "DefinitionProvider",
			shouldSupport: false,
		},
		{
			name:          "hover capability",
			capabilities:  []string{"hover"},
			checkFunc:     func(c *protocol.ServerCapabilities) bool { return c.HoverProvider != nil && c.HoverProvider.Value != nil },
			checkName:     "HoverProvider",
			shouldSupport: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := MockServerCapabilities{SupportedCapabilities: tt.capabilities}
			caps := mock.ToProtocol()

			result := tt.checkFunc(caps)
			if result != tt.shouldSupport {
				t.Errorf("%s support = %v, expected %v", tt.checkName, result, tt.shouldSupport)
			}
		})
	}
}
