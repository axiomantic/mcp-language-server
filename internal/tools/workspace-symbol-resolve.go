package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/isaacphi/mcp-language-server/internal/lsp"
	"github.com/isaacphi/mcp-language-server/internal/protocol"
)

// GetWorkspaceSymbolResolved searches for workspace symbols and resolves their details.
//
// This function performs two LSP operations:
// 1. workspace/symbol - searches for symbols matching the query
// 2. workspaceSymbol/resolve - resolves additional details for each symbol
//
// The resolve step is only performed if the server supports WorkspaceSymbolOptions.ResolveProvider.
// If resolve is not supported, returns the basic symbol information from workspace/symbol.
func GetWorkspaceSymbolResolved(ctx context.Context, client *lsp.Client, query string) (string, error) {
	// Step 1: Search for workspace symbols
	params := protocol.WorkspaceSymbolParams{
		Query: query,
	}

	result, err := client.Symbol(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to search workspace symbols: %v", err)
	}

	// Extract symbols from the Or_Result_workspace_symbol type
	var symbols []protocol.WorkspaceSymbol
	if result.Value != nil {
		switch v := result.Value.(type) {
		case []protocol.WorkspaceSymbol:
			symbols = v
		case []protocol.SymbolInformation:
			// Convert SymbolInformation to WorkspaceSymbol
			for _, si := range v {
				ws := protocol.WorkspaceSymbol{
					Location: protocol.Or_WorkspaceSymbol_location{Value: si.Location},
				}
				ws.Name = si.Name
				ws.Kind = si.Kind
				ws.ContainerName = si.ContainerName
				ws.Tags = si.Tags
				symbols = append(symbols, ws)
			}
		default:
			return "", fmt.Errorf("unexpected workspace symbol result type: %T", v)
		}
	}

	if len(symbols) == 0 {
		return fmt.Sprintf("No symbols found matching query: %s", query), nil
	}

	// Step 2: Resolve each symbol for additional details (if supported)
	caps := client.GetCapabilities()
	resolveSupported := lsp.HasWorkspaceSymbolResolveSupport(caps)

	var output strings.Builder
	output.WriteString(fmt.Sprintf("Found %d symbol(s) matching '%s':\n\n", len(symbols), query))

	for i, symbol := range symbols {
		output.WriteString(fmt.Sprintf("%d. %s (%s)\n", i+1, symbol.Name, symbolKindToString(symbol.Kind)))

		if symbol.ContainerName != "" {
			output.WriteString(fmt.Sprintf("   Container: %s\n", symbol.ContainerName))
		}

		// Try to resolve for more details if supported
		if resolveSupported {
			resolved, err := client.ResolveWorkspaceSymbol(ctx, symbol)
			if err == nil {
				// Use resolved symbol with additional details
				symbol = resolved
			}
			// Silently continue if resolve fails - we still have basic info
		}

		// Format location information
		if symbol.Location.Value != nil {
			switch loc := symbol.Location.Value.(type) {
			case protocol.Location:
				output.WriteString(fmt.Sprintf("   Location: %s:%d:%d\n",
					loc.URI,
					loc.Range.Start.Line+1,
					loc.Range.Start.Character+1))
			case protocol.LocationUriOnly:
				output.WriteString(fmt.Sprintf("   URI: %s\n", loc.URI))
			}
		}

		if symbol.Tags != nil && len(symbol.Tags) > 0 {
			output.WriteString(fmt.Sprintf("   Tags: %v\n", symbol.Tags))
		}

		output.WriteString("\n")
	}

	return output.String(), nil
}

// symbolKindToString converts a SymbolKind to its string representation
func symbolKindToString(kind protocol.SymbolKind) string {
	kinds := map[protocol.SymbolKind]string{
		1:  "File",
		2:  "Module",
		3:  "Namespace",
		4:  "Package",
		5:  "Class",
		6:  "Method",
		7:  "Property",
		8:  "Field",
		9:  "Constructor",
		10: "Enum",
		11: "Interface",
		12: "Function",
		13: "Variable",
		14: "Constant",
		15: "String",
		16: "Number",
		17: "Boolean",
		18: "Array",
		19: "Object",
		20: "Key",
		21: "Null",
		22: "EnumMember",
		23: "Struct",
		24: "Event",
		25: "Operator",
		26: "TypeParameter",
	}

	if s, ok := kinds[kind]; ok {
		return s
	}
	return fmt.Sprintf("Unknown(%d)", kind)
}
