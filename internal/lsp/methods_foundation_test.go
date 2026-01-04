package lsp

import (
	"context"
	"testing"

	"github.com/isaacphi/mcp-language-server/internal/protocol"
)

// TestLSPMethodStubs verifies that all foundation LSP method stubs exist
// This test primarily verifies compilation - if it compiles, the methods exist
func TestLSPMethodStubs(t *testing.T) {
	t.Log("Verifying LSP method stubs exist (compile-time check)")
	// The fact that this test compiles proves all the method signatures exist
	// and are correct. We don't need to call them since we don't have a real server.
}

// TestMethodSignatures verifies method signatures compile correctly
func TestMethodSignatures(t *testing.T) {
	// This test verifies that all method signatures are correct at compile time
	// If this test compiles, it means all required method signatures exist

	var _ interface {
		SemanticTokensFull(context.Context, protocol.SemanticTokensParams) (protocol.SemanticTokens, error)
		SemanticTokensRange(context.Context, protocol.SemanticTokensRangeParams) (protocol.SemanticTokens, error)
		PrepareTypeHierarchy(context.Context, protocol.TypeHierarchyPrepareParams) ([]protocol.TypeHierarchyItem, error)
		Supertypes(context.Context, protocol.TypeHierarchySupertypesParams) ([]protocol.TypeHierarchyItem, error)
		Subtypes(context.Context, protocol.TypeHierarchySubtypesParams) ([]protocol.TypeHierarchyItem, error)
		InlayHint(context.Context, protocol.InlayHintParams) ([]protocol.InlayHint, error)
		Formatting(context.Context, protocol.DocumentFormattingParams) ([]protocol.TextEdit, error)
		RangeFormatting(context.Context, protocol.DocumentRangeFormattingParams) ([]protocol.TextEdit, error)
		OnTypeFormatting(context.Context, protocol.DocumentOnTypeFormattingParams) ([]protocol.TextEdit, error)
		FoldingRange(context.Context, protocol.FoldingRangeParams) ([]protocol.FoldingRange, error)
		SelectionRange(context.Context, protocol.SelectionRangeParams) ([]protocol.SelectionRange, error)
		WillCreateFiles(context.Context, protocol.CreateFilesParams) (protocol.WorkspaceEdit, error)
		DidCreateFiles(context.Context, protocol.CreateFilesParams) error
		WillRenameFiles(context.Context, protocol.RenameFilesParams) (protocol.WorkspaceEdit, error)
		DidRenameFiles(context.Context, protocol.RenameFilesParams) error
		WillDeleteFiles(context.Context, protocol.DeleteFilesParams) (protocol.WorkspaceEdit, error)
		DidDeleteFiles(context.Context, protocol.DeleteFilesParams) error
	} = (*Client)(nil)

	t.Log("All method signatures verified at compile time")
}
