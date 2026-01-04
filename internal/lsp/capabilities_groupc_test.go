package lsp

import (
	"testing"

	"github.com/isaacphi/mcp-language-server/internal/protocol"
)

func TestHasFoldingRangeSupport(t *testing.T) {
	tests := []struct {
		name     string
		caps     *protocol.ServerCapabilities
		expected bool
	}{
		{
			name: "folding range supported",
			caps: &protocol.ServerCapabilities{
				FoldingRangeProvider: &protocol.Or_ServerCapabilities_foldingRangeProvider{
					Value: true,
				},
			},
			expected: true,
		},
		{
			name: "folding range Value nil",
			caps: &protocol.ServerCapabilities{
				FoldingRangeProvider: &protocol.Or_ServerCapabilities_foldingRangeProvider{
					Value: nil,
				},
			},
			expected: false,
		},
		{
			name: "folding range provider nil",
			caps: &protocol.ServerCapabilities{
				FoldingRangeProvider: nil,
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
			result := HasFoldingRangeSupport(tt.caps)
			if result != tt.expected {
				t.Errorf("HasFoldingRangeSupport() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestHasSelectionRangeSupport(t *testing.T) {
	tests := []struct {
		name     string
		caps     *protocol.ServerCapabilities
		expected bool
	}{
		{
			name: "selection range supported",
			caps: &protocol.ServerCapabilities{
				SelectionRangeProvider: &protocol.Or_ServerCapabilities_selectionRangeProvider{
					Value: true,
				},
			},
			expected: true,
		},
		{
			name: "selection range Value nil",
			caps: &protocol.ServerCapabilities{
				SelectionRangeProvider: &protocol.Or_ServerCapabilities_selectionRangeProvider{
					Value: nil,
				},
			},
			expected: false,
		},
		{
			name: "selection range provider nil",
			caps: &protocol.ServerCapabilities{
				SelectionRangeProvider: nil,
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
			result := HasSelectionRangeSupport(tt.caps)
			if result != tt.expected {
				t.Errorf("HasSelectionRangeSupport() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
