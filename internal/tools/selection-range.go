package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/isaacphi/mcp-language-server/internal/lsp"
	"github.com/isaacphi/mcp-language-server/internal/protocol"
)

// GetSelectionRanges retrieves selection ranges for a position in a document.
// Selection ranges show the hierarchy of nested selections at a given position,
// allowing progressive expansion of the selection.
func GetSelectionRanges(ctx context.Context, client *lsp.Client, filePath string, line, column int) (string, error) {
	// Open the file if not already open
	err := client.OpenFile(ctx, filePath)
	if err != nil {
		return "", fmt.Errorf("could not open file: %v", err)
	}

	// Convert 1-indexed line/column to 0-indexed for LSP protocol
	position := protocol.Position{
		Line:      uint32(line - 1),
		Character: uint32(column - 1),
	}

	params := protocol.SelectionRangeParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: protocol.DocumentUri("file://" + filePath),
		},
		Positions: []protocol.Position{position},
	}

	// Execute the selection range request
	selectionRanges, err := client.SelectionRange(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to get selection ranges: %v", err)
	}

	if len(selectionRanges) == 0 {
		return "No selection ranges found", nil
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("Selection Ranges for %s at line %d, column %d:\n\n",
		filePath, line, column))

	// Walk the selection hierarchy for the first (and only) position
	selRange := &selectionRanges[0]
	level := 0
	for selRange != nil {
		startLine := selRange.Range.Start.Line + 1
		startChar := selRange.Range.Start.Character + 1
		endLine := selRange.Range.End.Line + 1
		endChar := selRange.Range.End.Character + 1

		indent := strings.Repeat("  ", level)
		output.WriteString(fmt.Sprintf("%sLevel %d: [%d:%d-%d:%d]\n",
			indent, level, startLine, startChar, endLine, endChar))

		selRange = selRange.Parent
		level++
	}

	return output.String(), nil
}
