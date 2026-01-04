package lsp

import "github.com/isaacphi/mcp-language-server/internal/protocol"

// HasDefinitionSupport checks if the server supports textDocument/definition
// AND workspace/symbol (both required by our definition tool implementation).
//
// The definition tool uses workspace/symbol to locate symbols by name (step 1),
// then textDocument/definition to retrieve the actual source code (step 2).
// Verified in internal/tools/definition.go:13-17.
//
// CRITICAL: Uses two-part check for Or_* types (pointer != nil && .Value != nil).
// See design doc Section 2.3 for Or_* type behavior.
func HasDefinitionSupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	return caps.DefinitionProvider != nil &&
		caps.DefinitionProvider.Value != nil &&
		caps.WorkspaceSymbolProvider != nil &&
		caps.WorkspaceSymbolProvider.Value != nil
}

// HasReferencesSupport checks if the server supports textDocument/references.
//
// CRITICAL: Uses two-part check for Or_* type (pointer != nil && .Value != nil).
func HasReferencesSupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	return caps.ReferencesProvider != nil &&
		caps.ReferencesProvider.Value != nil
}

// HasHoverSupport checks if the server supports textDocument/hover.
//
// CRITICAL: Uses two-part check for Or_* type (pointer != nil && .Value != nil).
func HasHoverSupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	return caps.HoverProvider != nil &&
		caps.HoverProvider.Value != nil
}

// HasDocumentSymbolSupport checks if the server supports textDocument/documentSymbol.
//
// CRITICAL: Uses two-part check for Or_* type (pointer != nil && .Value != nil).
func HasDocumentSymbolSupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	return caps.DocumentSymbolProvider != nil &&
		caps.DocumentSymbolProvider.Value != nil
}

// HasCallHierarchySupport checks if the server supports call hierarchy
// (textDocument/prepareCallHierarchy, callHierarchy/incomingCalls, callHierarchy/outgoingCalls).
//
// Call Hierarchy was added in LSP 3.16.0.
//
// CRITICAL: Uses two-part check for Or_* type (pointer != nil && .Value != nil).
func HasCallHierarchySupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	return caps.CallHierarchyProvider != nil &&
		caps.CallHierarchyProvider.Value != nil
}

// HasWorkspaceSymbolSupport checks if the server supports workspace/symbol.
//
// Used by definition tool as a dependency check.
//
// CRITICAL: Uses two-part check for Or_* type (pointer != nil && .Value != nil).
func HasWorkspaceSymbolSupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	return caps.WorkspaceSymbolProvider != nil &&
		caps.WorkspaceSymbolProvider.Value != nil
}

// HasRenameSupport checks if the server supports textDocument/rename.
//
// RenameProvider is interface{} type - can be bool or RenameOptions.
// Simple nil check is sufficient (no .Value field to check).
func HasRenameSupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	return caps.RenameProvider != nil
}

// HasCodeActionSupport checks if the server supports textDocument/codeAction.
//
// CodeActionProvider is interface{} type - can be bool or CodeActionOptions.
// Simple nil check is sufficient (no .Value field to check).
func HasCodeActionSupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	return caps.CodeActionProvider != nil
}

// HasSignatureHelpSupport checks if the server supports textDocument/signatureHelp.
//
// SignatureHelpProvider is *SignatureHelpOptions type.
// Simple nil check is sufficient (pointer type, not Or_* type).
func HasSignatureHelpSupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	return caps.SignatureHelpProvider != nil
}

// HasCodeLensSupport checks if the server supports textDocument/codeLens.
//
// CodeLensProvider is *CodeLensOptions type.
// Simple nil check is sufficient (pointer type, not Or_* type).
func HasCodeLensSupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	return caps.CodeLensProvider != nil
}

// HasFoldingRangeSupport checks if the server supports textDocument/foldingRange.
//
// CRITICAL: Uses two-part check for Or_* type (pointer != nil && .Value != nil).
func HasFoldingRangeSupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	return caps.FoldingRangeProvider != nil &&
		caps.FoldingRangeProvider.Value != nil
}

// HasSelectionRangeSupport checks if the server supports textDocument/selectionRange.
//
// CRITICAL: Uses two-part check for Or_* type (pointer != nil && .Value != nil).
func HasSelectionRangeSupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	return caps.SelectionRangeProvider != nil &&
		caps.SelectionRangeProvider.Value != nil
}

// HasSemanticTokensSupport checks if the server supports textDocument/semanticTokens.
//
// SemanticTokensProvider is interface{} type - can be SemanticTokensOptions or SemanticTokensRegistrationOptions.
// Simple nil check is sufficient (no .Value field to check).
func HasSemanticTokensSupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	return caps.SemanticTokensProvider != nil
}

// HasTypeHierarchySupport checks if the server supports type hierarchy
// (textDocument/prepareTypeHierarchy, typeHierarchy/supertypes, typeHierarchy/subtypes).
//
// Type Hierarchy was added in LSP 3.17.0.
//
// CRITICAL: Uses two-part check for Or_* type (pointer != nil && .Value != nil).
func HasTypeHierarchySupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	return caps.TypeHierarchyProvider != nil &&
		caps.TypeHierarchyProvider.Value != nil
}

// HasInlayHintSupport checks if the server supports textDocument/inlayHint.
//
// InlayHintProvider is interface{} type - can be InlayHintOptions or InlayHintRegistrationOptions.
// Simple nil check is sufficient (no .Value field to check).
func HasInlayHintSupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	return caps.InlayHintProvider != nil
}

// AlwaysSupported returns true for core tools that don't require capability checks.
//
// Core tools:
// - edit_file: Requires TextDocumentSync, which every LSP server must provide
// - diagnostics: Uses push notifications (textDocument/publishDiagnostics), not capability-based
//
// This function exists for documentation and consistency with the capability check pattern.
func AlwaysSupported() bool {
	return true
}

// HasWorkspaceSymbolResolveSupport checks if the server supports workspace/resolveWorkspaceSymbol.
//
// This capability is part of WorkspaceSymbolOptions.ResolveProvider (LSP 3.17+).
// The server can provide workspace symbols and then resolve additional details for them.
//
// CRITICAL: Uses two-part check for Or_* type (pointer != nil && .Value != nil),
// then checks if Value is WorkspaceSymbolOptions with ResolveProvider = true.
func HasWorkspaceSymbolResolveSupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	if caps.WorkspaceSymbolProvider == nil || caps.WorkspaceSymbolProvider.Value == nil {
		return false
	}

	// Check if Value is WorkspaceSymbolOptions with ResolveProvider set
	if opts, ok := caps.WorkspaceSymbolProvider.Value.(protocol.WorkspaceSymbolOptions); ok {
		return opts.ResolveProvider
	}

	// If it's just a bool (not WorkspaceSymbolOptions), no resolve support
	return false
}

// HasPrepareRenameSupport checks if the server supports textDocument/prepareRename.
//
// PrepareRename allows validating a rename operation before executing it.
// This is indicated by RenameProvider being a struct with prepareProvider field.
//
// RenameProvider is interface{} type - can be bool or map[string]interface{}.
func HasPrepareRenameSupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	if caps.RenameProvider == nil {
		return false
	}

	// Check if it's a map with prepareProvider field
	if opts, ok := caps.RenameProvider.(map[string]interface{}); ok {
		if prepareProvider, exists := opts["prepareProvider"]; exists {
			if prepare, ok := prepareProvider.(bool); ok {
				return prepare
			}
		}
	}

	return false
}

// HasFormattingSupport checks if the server supports textDocument/formatting.
//
// CRITICAL: Uses two-part check for Or_* type (pointer != nil && .Value != nil).
func HasFormattingSupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	return caps.DocumentFormattingProvider != nil &&
		caps.DocumentFormattingProvider.Value != nil
}

// HasRangeFormattingSupport checks if the server supports textDocument/rangeFormatting.
//
// CRITICAL: Uses two-part check for Or_* type (pointer != nil && .Value != nil).
func HasRangeFormattingSupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	return caps.DocumentRangeFormattingProvider != nil &&
		caps.DocumentRangeFormattingProvider.Value != nil
}

// HasOnTypeFormattingSupport checks if the server supports textDocument/onTypeFormatting.
//
// DocumentOnTypeFormattingProvider is *DocumentOnTypeFormattingOptions type.
// Simple nil check is sufficient (pointer type, not Or_* type).
func HasOnTypeFormattingSupport(caps *protocol.ServerCapabilities) bool {
	if caps == nil {
		return false
	}
	return caps.DocumentOnTypeFormattingProvider != nil
}
