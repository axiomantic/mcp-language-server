package tools

import (
	"context"
	"testing"

	"github.com/isaacphi/mcp-language-server/internal/lsp"
)

func TestGetInlayHints(t *testing.T) {
	t.Run("returns error when file cannot be opened", func(t *testing.T) {
		ctx := context.Background()
		client := &lsp.Client{} // This will fail to open file

		_, err := GetInlayHints(ctx, client, "/nonexistent/file.go", 1, 10)
		if err == nil {
			t.Error("Expected error when opening nonexistent file, got nil")
		}
	})
}
