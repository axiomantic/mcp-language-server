package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/isaacphi/mcp-language-server/internal/lsp"
	"github.com/isaacphi/mcp-language-server/internal/protocol"
)

// GetInlayHints retrieves inlay hints for a range of lines in a file
// Inlay hints provide inline annotations like parameter names and type hints
func GetInlayHints(ctx context.Context, client *lsp.Client, filePath string, startLine, endLine int) (string, error) {
	// Open the file if not already open
	err := client.OpenFile(ctx, filePath)
	if err != nil {
		return "", fmt.Errorf("could not open file: %v", err)
	}

	params := protocol.InlayHintParams{}
	uri := protocol.DocumentUri("file://" + filePath)
	params.TextDocument = protocol.TextDocumentIdentifier{
		URI: uri,
	}

	// Set the range for inlay hints (convert 1-indexed to 0-indexed)
	params.Range = protocol.Range{
		Start: protocol.Position{
			Line:      uint32(startLine - 1),
			Character: 0,
		},
		End: protocol.Position{
			Line:      uint32(endLine),
			Character: 0,
		},
	}

	// Execute the inlay hint request
	hints, err := client.InlayHint(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to get inlay hints: %v", err)
	}

	if len(hints) == 0 {
		return fmt.Sprintf("No inlay hints found for lines %d-%d", startLine, endLine), nil
	}

	var result strings.Builder

	result.WriteString(fmt.Sprintf("Inlay Hints for lines %d-%d (%d hints):\n\n", startLine, endLine, len(hints)))

	for _, hint := range hints {
		line := hint.Position.Line + 1
		char := hint.Position.Character + 1

		// Extract label text from the hint
		var labelText string
		if len(hint.Label) > 0 {
			// Label is an array of InlayHintLabelPart
			var parts []string
			for _, part := range hint.Label {
				parts = append(parts, part.Value)
			}
			labelText = strings.Join(parts, "")
		}

		// Determine hint kind
		kindStr := "hint"
		if hint.Kind == 1 {
			kindStr = "type"
		} else if hint.Kind == 2 {
			kindStr = "parameter"
		}

		result.WriteString(fmt.Sprintf("Line %d, Col %d [%s]: %s",
			line, char, kindStr, labelText))

		// Add tooltip if available
		if hint.Tooltip != nil {
			tooltipText := fmt.Sprintf("%v", hint.Tooltip.Value)
			// Truncate long tooltips
			if len(tooltipText) > 100 {
				tooltipText = tooltipText[:97] + "..."
			}
			if tooltipText != "" && tooltipText != "<nil>" {
				result.WriteString(fmt.Sprintf(" (tooltip: %s)", tooltipText))
			}
		}

		result.WriteString("\n")
	}

	return result.String(), nil
}
