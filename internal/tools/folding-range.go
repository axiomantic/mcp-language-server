package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/isaacphi/mcp-language-server/internal/lsp"
	"github.com/isaacphi/mcp-language-server/internal/protocol"
)

// GetFoldingRanges retrieves folding ranges for a document
func GetFoldingRanges(ctx context.Context, client *lsp.Client, filePath string) (string, error) {
	// Open the file if not already open
	err := client.OpenFile(ctx, filePath)
	if err != nil {
		return "", fmt.Errorf("could not open file: %v", err)
	}

	params := protocol.FoldingRangeParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: protocol.DocumentUri("file://" + filePath),
		},
	}

	// Execute the folding range request
	foldingRanges, err := client.FoldingRange(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to get folding ranges: %v", err)
	}

	if len(foldingRanges) == 0 {
		return "No folding ranges found", nil
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("Folding Ranges for %s:\n\n", filePath))

	for i, fr := range foldingRanges {
		// Convert from 0-indexed to 1-indexed for display
		startLine := fr.StartLine + 1
		endLine := fr.EndLine + 1

		kind := fr.Kind
		if kind == "" {
			kind = "region"
		}

		output.WriteString(fmt.Sprintf("%d. %s: lines %d-%d",
			i+1, kind, startLine, endLine))

		// Add character positions if specified
		if fr.StartCharacter > 0 || fr.EndCharacter > 0 {
			output.WriteString(fmt.Sprintf(" (chars %d-%d)",
				fr.StartCharacter+1, fr.EndCharacter+1))
		}

		output.WriteString("\n")
	}

	return output.String(), nil
}
