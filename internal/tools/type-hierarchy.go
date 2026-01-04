package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/isaacphi/mcp-language-server/internal/lsp"
	"github.com/isaacphi/mcp-language-server/internal/protocol"
)

// GetTypeHierarchy retrieves type hierarchy (supertypes, subtypes, or both) for a type at the specified position
func GetTypeHierarchy(ctx context.Context, client *lsp.Client, filePath string, line, column int, direction string) (string, error) {
	// Validate direction parameter
	if direction != "supertypes" && direction != "subtypes" && direction != "both" {
		return "", fmt.Errorf("direction must be 'supertypes', 'subtypes', or 'both'")
	}

	// Open the file if not already open
	err := client.OpenFile(ctx, filePath)
	if err != nil {
		return "", fmt.Errorf("could not open file: %v", err)
	}

	// Prepare type hierarchy (find the type at the given position)
	prepareParams := protocol.TypeHierarchyPrepareParams{}
	position := protocol.Position{
		Line:      uint32(line - 1),
		Character: uint32(column - 1),
	}
	uri := protocol.DocumentUri("file://" + filePath)
	prepareParams.TextDocument = protocol.TextDocumentIdentifier{
		URI: uri,
	}
	prepareParams.Position = position

	items, err := client.PrepareTypeHierarchy(ctx, prepareParams)
	if err != nil {
		return "", fmt.Errorf("failed to prepare type hierarchy: %v", err)
	}

	if len(items) == 0 {
		return "No type found at this position", nil
	}

	var result strings.Builder

	// Process each type hierarchy item
	for _, item := range items {
		result.WriteString(fmt.Sprintf("Type: %s\n", item.Name))
		result.WriteString(fmt.Sprintf("Kind: %d\n", item.Kind))
		if item.Detail != "" {
			result.WriteString(fmt.Sprintf("Detail: %s\n", item.Detail))
		}
		result.WriteString(fmt.Sprintf("Location: %s:%d:%d\n\n",
			item.URI,
			item.Range.Start.Line+1,
			item.Range.Start.Character+1))

		// Get supertypes if requested
		if direction == "supertypes" || direction == "both" {
			supertypesParams := protocol.TypeHierarchySupertypesParams{
				Item: item,
			}
			supertypes, err := client.Supertypes(ctx, supertypesParams)
			if err != nil {
				toolsLogger.Warn("failed to get supertypes: %v", err)
			} else if len(supertypes) > 0 {
				result.WriteString("Supertypes:\n")
				for _, supertype := range supertypes {
					result.WriteString(fmt.Sprintf("  - %s (%d)\n", supertype.Name, supertype.Kind))
					if supertype.Detail != "" {
						result.WriteString(fmt.Sprintf("    Detail: %s\n", supertype.Detail))
					}
				}
				result.WriteString("\n")
			}
		}

		// Get subtypes if requested
		if direction == "subtypes" || direction == "both" {
			subtypesParams := protocol.TypeHierarchySubtypesParams{
				Item: item,
			}
			subtypes, err := client.Subtypes(ctx, subtypesParams)
			if err != nil {
				toolsLogger.Warn("failed to get subtypes: %v", err)
			} else if len(subtypes) > 0 {
				result.WriteString("Subtypes:\n")
				for _, subtype := range subtypes {
					result.WriteString(fmt.Sprintf("  - %s (%d)\n", subtype.Name, subtype.Kind))
					if subtype.Detail != "" {
						result.WriteString(fmt.Sprintf("    Detail: %s\n", subtype.Detail))
					}
				}
				result.WriteString("\n")
			}
		}
	}

	return result.String(), nil
}
