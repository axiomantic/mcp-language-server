package tools

import (
	"context"
	"testing"

	"github.com/isaacphi/mcp-language-server/internal/lsp"
)

// TestGetWorkspaceSymbolResolved_Signature tests the function signature
func TestGetWorkspaceSymbolResolved_Signature(t *testing.T) {
	// Verify that GetWorkspaceSymbolResolved function exists and has correct signature
	var _ func(context.Context, *lsp.Client, string) (string, error) = GetWorkspaceSymbolResolved
}
