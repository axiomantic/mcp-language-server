package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/isaacphi/mcp-language-server/internal/lsp"
	"github.com/isaacphi/mcp-language-server/internal/protocol"
)

// GetSemanticTokens retrieves semantic tokens for a file, providing syntax highlighting information
func GetSemanticTokens(ctx context.Context, client *lsp.Client, filePath string) (string, error) {
	// Open the file if not already open
	err := client.OpenFile(ctx, filePath)
	if err != nil {
		return "", fmt.Errorf("could not open file: %v", err)
	}

	params := protocol.SemanticTokensParams{}
	uri := protocol.DocumentUri("file://" + filePath)
	params.TextDocument = protocol.TextDocumentIdentifier{
		URI: uri,
	}

	// Execute the semantic tokens request
	tokensResult, err := client.SemanticTokensFull(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to get semantic tokens: %v", err)
	}

	var result strings.Builder

	// Get the legend from server capabilities to decode token types and modifiers
	caps := client.GetCapabilities()
	if caps == nil || caps.SemanticTokensProvider == nil {
		return "", fmt.Errorf("server does not support semantic tokens")
	}

	// Extract legend from the semantic tokens provider
	var tokenTypes []string
	var tokenModifiers []string

	// SemanticTokensProvider can be SemanticTokensOptions or SemanticTokensRegistrationOptions
	// Both have a Legend field
	if opts, ok := caps.SemanticTokensProvider.(map[string]interface{}); ok {
		if legend, ok := opts["legend"].(map[string]interface{}); ok {
			if types, ok := legend["tokenTypes"].([]interface{}); ok {
				for _, t := range types {
					if ts, ok := t.(string); ok {
						tokenTypes = append(tokenTypes, ts)
					}
				}
			}
			if mods, ok := legend["tokenModifiers"].([]interface{}); ok {
				for _, m := range mods {
					if ms, ok := m.(string); ok {
						tokenModifiers = append(tokenModifiers, ms)
					}
				}
			}
		}
	}

	if len(tokenTypes) == 0 {
		return "", fmt.Errorf("no token types found in server capabilities")
	}

	// Decode the semantic tokens (delta-encoded integers)
	tokens := tokensResult.Data
	if len(tokens) == 0 {
		return "No semantic tokens found in this file", nil
	}

	result.WriteString(fmt.Sprintf("Semantic Tokens (%d tokens):\n\n", len(tokens)/5))
	result.WriteString("Token Types Available:\n")
	for i, tokenType := range tokenTypes {
		result.WriteString(fmt.Sprintf("  %d: %s\n", i, tokenType))
	}
	result.WriteString("\nToken Modifiers Available:\n")
	for i, tokenMod := range tokenModifiers {
		result.WriteString(fmt.Sprintf("  %d: %s\n", i, tokenMod))
	}

	result.WriteString("\nTokens (first 10):\n")
	line := uint32(0)
	char := uint32(0)

	// Decode tokens (each token is 5 integers: deltaLine, deltaChar, length, tokenType, tokenModifiers)
	for i := 0; i < len(tokens) && i/5 < 10; i += 5 {
		deltaLine := tokens[i]
		deltaChar := tokens[i+1]
		length := tokens[i+2]
		tokenType := tokens[i+3]
		tokenMods := tokens[i+4]

		// Update position
		line += deltaLine
		if deltaLine > 0 {
			char = 0
		}
		char += deltaChar

		// Get token type name
		tokenTypeName := "unknown"
		if int(tokenType) < len(tokenTypes) {
			tokenTypeName = tokenTypes[tokenType]
		}

		// Decode token modifiers (bit flags)
		var modNames []string
		for j := 0; j < len(tokenModifiers); j++ {
			if tokenMods&(1<<uint(j)) != 0 {
				modNames = append(modNames, tokenModifiers[j])
			}
		}

		modStr := "none"
		if len(modNames) > 0 {
			modStr = strings.Join(modNames, ", ")
		}

		result.WriteString(fmt.Sprintf("  Line %d, Col %d, Len %d: %s [%s]\n",
			line+1, char+1, length, tokenTypeName, modStr))
	}

	if len(tokens)/5 > 10 {
		result.WriteString(fmt.Sprintf("\n... and %d more tokens\n", len(tokens)/5-10))
	}

	return result.String(), nil
}
