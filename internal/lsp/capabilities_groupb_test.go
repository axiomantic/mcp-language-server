package lsp

import (
	"testing"

	"github.com/isaacphi/mcp-language-server/internal/protocol"
)

func TestHasWorkspaceSymbolResolveSupport(t *testing.T) {
	tests := []struct {
		name     string
		caps     *protocol.ServerCapabilities
		expected bool
	}{
		{
			name: "workspace symbol resolve supported",
			caps: &protocol.ServerCapabilities{
				WorkspaceSymbolProvider: &protocol.Or_ServerCapabilities_workspaceSymbolProvider{
					Value: protocol.WorkspaceSymbolOptions{
						ResolveProvider: true,
					},
				},
			},
			expected: true,
		},
		{
			name: "workspace symbol without resolve",
			caps: &protocol.ServerCapabilities{
				WorkspaceSymbolProvider: &protocol.Or_ServerCapabilities_workspaceSymbolProvider{
					Value: protocol.WorkspaceSymbolOptions{
						ResolveProvider: false,
					},
				},
			},
			expected: false,
		},
		{
			name: "workspace symbol as bool (no resolve)",
			caps: &protocol.ServerCapabilities{
				WorkspaceSymbolProvider: &protocol.Or_ServerCapabilities_workspaceSymbolProvider{
					Value: true,
				},
			},
			expected: false,
		},
		{
			name: "workspace symbol provider nil",
			caps: &protocol.ServerCapabilities{
				WorkspaceSymbolProvider: nil,
			},
			expected: false,
		},
		{
			name:     "nil capabilities",
			caps:     nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasWorkspaceSymbolResolveSupport(tt.caps)
			if result != tt.expected {
				t.Errorf("HasWorkspaceSymbolResolveSupport() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestHasPrepareRenameSupport(t *testing.T) {
	tests := []struct {
		name     string
		caps     *protocol.ServerCapabilities
		expected bool
	}{
		{
			name: "prepare rename supported (options struct)",
			caps: &protocol.ServerCapabilities{
				RenameProvider: map[string]interface{}{"prepareProvider": true},
			},
			expected: true,
		},
		{
			name: "prepare rename explicitly false",
			caps: &protocol.ServerCapabilities{
				RenameProvider: map[string]interface{}{"prepareProvider": false},
			},
			expected: false,
		},
		{
			name: "rename supported as bool (no prepare)",
			caps: &protocol.ServerCapabilities{
				RenameProvider: true,
			},
			expected: false,
		},
		{
			name: "rename provider nil",
			caps: &protocol.ServerCapabilities{
				RenameProvider: nil,
			},
			expected: false,
		},
		{
			name:     "nil capabilities",
			caps:     nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasPrepareRenameSupport(tt.caps)
			if result != tt.expected {
				t.Errorf("HasPrepareRenameSupport() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestHasFormattingSupport(t *testing.T) {
	tests := []struct {
		name     string
		caps     *protocol.ServerCapabilities
		expected bool
	}{
		{
			name: "formatting supported",
			caps: &protocol.ServerCapabilities{
				DocumentFormattingProvider: &protocol.Or_ServerCapabilities_documentFormattingProvider{
					Value: true,
				},
			},
			expected: true,
		},
		{
			name: "formatting Value nil",
			caps: &protocol.ServerCapabilities{
				DocumentFormattingProvider: &protocol.Or_ServerCapabilities_documentFormattingProvider{
					Value: nil,
				},
			},
			expected: false,
		},
		{
			name: "formatting provider nil",
			caps: &protocol.ServerCapabilities{
				DocumentFormattingProvider: nil,
			},
			expected: false,
		},
		{
			name:     "nil capabilities",
			caps:     nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasFormattingSupport(tt.caps)
			if result != tt.expected {
				t.Errorf("HasFormattingSupport() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestHasRangeFormattingSupport(t *testing.T) {
	tests := []struct {
		name     string
		caps     *protocol.ServerCapabilities
		expected bool
	}{
		{
			name: "range formatting supported",
			caps: &protocol.ServerCapabilities{
				DocumentRangeFormattingProvider: &protocol.Or_ServerCapabilities_documentRangeFormattingProvider{
					Value: true,
				},
			},
			expected: true,
		},
		{
			name: "range formatting Value nil",
			caps: &protocol.ServerCapabilities{
				DocumentRangeFormattingProvider: &protocol.Or_ServerCapabilities_documentRangeFormattingProvider{
					Value: nil,
				},
			},
			expected: false,
		},
		{
			name: "range formatting provider nil",
			caps: &protocol.ServerCapabilities{
				DocumentRangeFormattingProvider: nil,
			},
			expected: false,
		},
		{
			name:     "nil capabilities",
			caps:     nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasRangeFormattingSupport(tt.caps)
			if result != tt.expected {
				t.Errorf("HasRangeFormattingSupport() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestHasOnTypeFormattingSupport(t *testing.T) {
	tests := []struct {
		name     string
		caps     *protocol.ServerCapabilities
		expected bool
	}{
		{
			name: "ontype formatting supported",
			caps: &protocol.ServerCapabilities{
				DocumentOnTypeFormattingProvider: &protocol.DocumentOnTypeFormattingOptions{
					FirstTriggerCharacter: ";",
				},
			},
			expected: true,
		},
		{
			name: "ontype formatting provider nil",
			caps: &protocol.ServerCapabilities{
				DocumentOnTypeFormattingProvider: nil,
			},
			expected: false,
		},
		{
			name:     "nil capabilities",
			caps:     nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasOnTypeFormattingSupport(tt.caps)
			if result != tt.expected {
				t.Errorf("HasOnTypeFormattingSupport() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
