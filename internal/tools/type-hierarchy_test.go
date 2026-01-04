package tools

import (
	"context"
	"testing"

	"github.com/isaacphi/mcp-language-server/internal/lsp"
)

func TestGetTypeHierarchy(t *testing.T) {
	t.Run("returns error when file cannot be opened", func(t *testing.T) {
		ctx := context.Background()
		client := &lsp.Client{} // This will fail to open file

		_, err := GetTypeHierarchy(ctx, client, "/nonexistent/file.go", 1, 1, "both")
		if err == nil {
			t.Error("Expected error when opening nonexistent file, got nil")
		}
	})

	t.Run("returns error for invalid direction", func(t *testing.T) {
		ctx := context.Background()
		client := &lsp.Client{}

		_, err := GetTypeHierarchy(ctx, client, "/some/file.go", 1, 1, "invalid")
		if err == nil {
			t.Error("Expected error for invalid direction, got nil")
		}
		if err.Error() != "direction must be 'supertypes', 'subtypes', or 'both'" {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})
}
