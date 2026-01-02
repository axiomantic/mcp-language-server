#!/bin/bash
# Test script to verify capability reporting across different LSP servers

set -euo pipefail

echo "Testing LSP server capabilities..."
echo ""

# Test with gopls
echo "=== gopls (Go) ==="
timeout 5s go run . -workspace ./integrationtests/workspaces/go -lsp gopls 2>&1 | \
    grep -A 10 "=== LSP Server Capabilities ===" || echo "gopls test failed or timed out"
echo ""

# Test with typescript-language-server
echo "=== typescript-language-server ==="
timeout 5s go run . -workspace ./integrationtests/workspaces/typescript \
    -lsp typescript-language-server -- --stdio 2>&1 | \
    grep -A 10 "=== LSP Server Capabilities ===" || echo "typescript-language-server test failed or timed out"
echo ""

# Test with rust-analyzer
echo "=== rust-analyzer ==="
timeout 5s go run . -workspace ./integrationtests/workspaces/rust -lsp rust-analyzer 2>&1 | \
    grep -A 10 "=== LSP Server Capabilities ===" || echo "rust-analyzer test failed or timed out"
echo ""

# Test with clangd
echo "=== clangd (C++) ==="
timeout 5s go run . -workspace ./integrationtests/workspaces/cpp -lsp clangd 2>&1 | \
    grep -A 10 "=== LSP Server Capabilities ===" || echo "clangd test failed or timed out"
echo ""

# Test with pyright
echo "=== pyright (Python) ==="
timeout 5s go run . -workspace ./integrationtests/workspaces/python -lsp pyright-langserver -- --stdio 2>&1 | \
    grep -A 10 "=== LSP Server Capabilities ===" || echo "pyright-langserver test failed or timed out"
echo ""

echo "Capability testing complete"
