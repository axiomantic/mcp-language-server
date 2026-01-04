package tools

import (
	"context"
	"testing"

	"github.com/isaacphi/mcp-language-server/internal/lsp"
)

// TestFormatDocument_Signature tests the function signature
func TestFormatDocument_Signature(t *testing.T) {
	// Verify that FormatDocument function exists and has correct signature
	// mode: "full", "range", "ontype"
	var _ func(context.Context, *lsp.Client, string, string, int, int, int, int, string) (string, error) = FormatDocument
}
