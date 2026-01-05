# MCP Language Server

[![Go Tests](https://github.com/isaacphi/mcp-language-server/actions/workflows/go.yml/badge.svg)](https://github.com/isaacphi/mcp-language-server/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/isaacphi/mcp-language-server)](https://goreportcard.com/report/github.com/isaacphi/mcp-language-server)
[![GoDoc](https://pkg.go.dev/badge/github.com/isaacphi/mcp-language-server)](https://pkg.go.dev/github.com/isaacphi/mcp-language-server)
[![Go Version](https://img.shields.io/github/go-mod/go-version/isaacphi/mcp-language-server)](https://github.com/isaacphi/mcp-language-server/blob/main/go.mod)

This is an [MCP](https://modelcontextprotocol.io/introduction) server that runs and exposes a [language server](https://microsoft.github.io/language-server-protocol/) to LLMs. Not a language server for MCP, whatever that would be.

## Demo

`mcp-language-server` helps MCP enabled clients navigate codebases more easily by giving them access semantic tools like get definition, references, rename, and diagnostics.

![Demo](demo.gif)

## Setup

1. **Install Go**: Follow instructions at <https://golang.org/doc/install>
2. **Install or update this server**: `go install github.com/isaacphi/mcp-language-server@latest`
3. **Install a language server**: _follow one of the guides below_
4. **Configure your MCP client**: _follow one of the guides below_

## Transport Options

The MCP Language Server supports two transport modes:

### Stdio Transport (Default)

Standard input/output transport for direct process communication. This is the default mode and requires no additional flags.

```bash
mcp-language-server --workspace=/path/to/project --lsp=gopls
```

### HTTP Transport

HTTP transport for network-based communication, useful for web UIs or multiple clients.

```bash
mcp-language-server --workspace=/path/to/project --lsp=gopls --transport=http --port=8080
```

**Endpoint:** `http://localhost:8080/mcp/v1`

**CLI Flags:**
- `--transport` - Transport type: `stdio` (default) or `http`
- `--port` - Port for HTTP transport (default: 8080)

**Security Notice:** HTTP transport is designed for local development only. The server binds to localhost and does not include authentication. Do NOT expose the HTTP port to untrusted networks.

## File Operations

When files are created, renamed, or deleted in your workspace, the server sends `notifications/resources/updated` to all connected MCP clients. This allows clients to stay synchronized with workspace changes.

### Notification Format

**Method:** `notifications/resources/updated`

**Payload:**
```json
{
  "type": "created|renamed|deleted",
  "uris": ["file:///workspace/path/to/file.go"]
}
```

### Examples

**File Created:**
```json
{
  "type": "created",
  "uris": ["file:///workspace/main.go"]
}
```

**File Renamed:**
```json
{
  "type": "renamed",
  "uris": [
    "file:///workspace/old.go",
    "file:///workspace/new.go"
  ]
}
```

**File Deleted:**
```json
{
  "type": "deleted",
  "uris": ["file:///workspace/temp.go"]
}
```

### Client-Side Handling

MCP clients can subscribe to these notifications to invalidate caches or update UI:

```javascript
mcpClient.on('notifications/resources/updated', (params) => {
  console.log(`Files ${params.type}:`, params.uris);
  // Refresh your file cache/UI
});
```

**Note:** The server uses a 100ms debounce window to detect rename operations (delete + create pairs). Notifications are best-effort and clients should handle missed notifications gracefully.

<details>
  <summary>Go (gopls)</summary>
  <div>
    <p><strong>Install gopls</strong>: <code>go install golang.org/x/tools/gopls@latest</code></p>
    <p><strong>Configure your MCP client</strong>: This will be different but similar for each client. For Claude Desktop, add the following to <code>~/Library/Application\ Support/Claude/claude_desktop_config.json</code></p>

<pre>
{
  "mcpServers": {
    "language-server": {
      "command": "mcp-language-server",
      "args": ["--workspace", "/Users/you/dev/yourproject/", "--lsp", "gopls"],
      "env": {
        "PATH": "/opt/homebrew/bin:/Users/you/go/bin",
        "GOPATH": "/users/you/go",
        "GOCACHE": "/users/you/Library/Caches/go-build",
        "GOMODCACHE": "/Users/you/go/pkg/mod"
      }
    }
  }
}
</pre>

<p><strong>Note</strong>: Not all clients will need these environment variables. For Claude Desktop you will need to update the environment variables above based on your machine and username:</p>
<ul>
  <li><code>PATH</code> needs to contain the path to <code>go</code> and to <code>gopls</code>. Get this with <code>echo $(which go):$(which gopls)</code></li>
  <li><code>GOPATH</code>, <code>GOCACHE</code>, and <code>GOMODCACHE</code> may be different on your machine. These are the defaults.</li>
</ul>

  </div>
</details>
<details>
  <summary>Rust (rust-analyzer)</summary>
  <div>
    <p><strong>Install rust-analyzer</strong>: <code>rustup component add rust-analyzer</code></p>
    <p><strong>Configure your MCP client</strong>: This will be different but similar for each client. For Claude Desktop, add the following to <code>~/Library/Application\ Support/Claude/claude_desktop_config.json</code></p>

<pre>
{
  "mcpServers": {
    "language-server": {
      "command": "mcp-language-server",
      "args": [
        "--workspace",
        "/Users/you/dev/yourproject/",
        "--lsp",
        "rust-analyzer"
      ]
    }
  }
}
</pre>
  </div>
</details>
<details>
  <summary>Python (pyright)</summary>
  <div>
    <p><strong>Install pyright</strong>: <code>npm install -g pyright</code></p>
    <p><strong>Configure your MCP client</strong>: This will be different but similar for each client. For Claude Desktop, add the following to <code>~/Library/Application\ Support/Claude/claude_desktop_config.json</code></p>

<pre>
{
  "mcpServers": {
    "language-server": {
      "command": "mcp-language-server",
      "args": [
        "--workspace",
        "/Users/you/dev/yourproject/",
        "--lsp",
        "pyright-langserver",
        "--",
        "--stdio"
      ]
    }
  }
}
</pre>
  </div>
</details>
<details>
  <summary>Typescript (typescript-language-server)</summary>
  <div>
    <p><strong>Install typescript-language-server</strong>: <code>npm install -g typescript typescript-language-server</code></p>
    <p><strong>Configure your MCP client</strong>: This will be different but similar for each client. For Claude Desktop, add the following to <code>~/Library/Application\ Support/Claude/claude_desktop_config.json</code></p>

<pre>
{
  "mcpServers": {
    "language-server": {
      "command": "mcp-language-server",
      "args": [
        "--workspace",
        "/Users/you/dev/yourproject/",
        "--lsp",
        "typescript-language-server",
        "--",
        "--stdio"
      ]
    }
  }
}
</pre>
  </div>
</details>
<details>
  <summary>C/C++ (clangd)</summary>
  <div>
    <p><strong>Install clangd</strong>: Download prebuilt binaries from the <a href="https://github.com/clangd/clangd/releases">official LLVM releases page</a> or install via your system's package manager (e.g., <code>apt install clangd</code>, <code>brew install clangd</code>).</p>
    <p><strong>Configure your MCP client</strong>: This will be different but similar for each client. For Claude Desktop, add the following to <code>~/Library/Application\\ Support/Claude/claude_desktop_config.json</code></p>

<pre>
{
  "mcpServers": {
    "language-server": {
      "command": "mcp-language-server",
      "args": [
        "--workspace",
        "/Users/you/dev/yourproject/",
        "--lsp",
        "/path/to/your/clangd_binary",
        "--",
        "--compile-commands-dir=/path/to/yourproject/build_or_compile_commands_dir"
      ]
    }
  }
}
</pre>
    <p><strong>Note</strong>:</p>
    <ul>
      <li>Replace <code>/path/to/your/clangd_binary</code> with the actual path to your clangd executable.</li>
      <li><code>--compile-commands-dir</code> should point to the directory containing your <code>compile_commands.json</code> file (e.g., <code>./build</code>, <code>./cmake-build-debug</code>).</li>
      <li>Ensure <code>compile_commands.json</code> is generated for your project for clangd to work effectively.</li>
    </ul>
  </div>
</details>
<details>
  <summary>Other</summary>
  <div>
    <p>I have only tested this repo with the servers above but it should be compatible with many more. Note:</p>
    <ul>
      <li>The language server must communicate over stdio.</li>
      <li>Any aruments after <code>--</code> are sent as arguments to the language server.</li>
      <li>Any env variables are passed on to the language server.</li>
    </ul>
  </div>
</details>

## Tool Availability

The MCP Language Server **dynamically registers tools** based on the capabilities advertised by the underlying LSP server. Not all tools may be available for all language servers.

### Core Tools (Always Available)

These tools are always registered regardless of LSP server capabilities:

- **`edit_file`** - Apply text edits to files (requires `TextDocumentSync`, which all LSP servers provide)
- **`diagnostics`** - Get diagnostic information (uses push notifications, not capability-based)

### Capability-Dependent Tools

The following tools are only available if the LSP server supports them:

- **`definition`** - Find symbol definitions
  - Requires: `DefinitionProvider` + `WorkspaceSymbolProvider`
  - Why both: Uses workspace/symbol to locate symbols, then definition to get code

- **`references`** - Find all symbol references
  - Requires: `ReferencesProvider`

- **`hover`** - Get hover information (types, documentation)
  - Requires: `HoverProvider`

- **`rename_symbol`** - Rename symbols across the codebase
  - Requires: `RenameProvider`

- **`code_actions`** - Get available quick fixes and refactorings
  - Requires: `CodeActionProvider`

- **`signature_help`** - Get function/method signature information
  - Requires: `SignatureHelpProvider`

- **`document_symbols`** - Get hierarchical symbol outline
  - Requires: `DocumentSymbolProvider`

- **`call_hierarchy`** - Find callers/callees of functions
  - Requires: `CallHierarchyProvider` (LSP 3.16+)

- **`get_codelens`** - Get code lens hints
  - Requires: `CodeLensProvider`

- **`execute_codelens`** - Execute code lens commands
  - Requires: `CodeLensProvider`

### Checking Available Tools

When starting the server, check the logs for capability information:

```
INFO: === LSP Server Capabilities ===
INFO: Definition: true
INFO: References: true
INFO: Hover: true
INFO: Rename: true
INFO: Code Actions: true
INFO: Code Lens: false
INFO: Signature Help: true
INFO: Document Symbols: true
INFO: Call Hierarchy: true
INFO: Workspace Symbols: true
INFO: ===============================
INFO: Registering core tools
DEBUG: Registering 'definition' tool
DEBUG: Registering 'references' tool
...
INFO: Skipping 'get_codelens' and 'execute_codelens' tools - LSP server doesn't support CodeLens capability
```

### Language Server Support Matrix

This matrix documents observed capability support across common LSP servers:

| Tool | gopls | typescript-language-server | rust-analyzer | clangd | pyright |
|------|-------|---------------------------|---------------|---------|---------|
| edit_file | ✅ | ✅ | ✅ | ✅ | ✅ |
| diagnostics | ✅ | ✅ | ✅ | ✅ | ✅ |
| definition | ✅ | ✅ | ✅ | ✅ | ✅ |
| references | ✅ | ✅ | ✅ | ✅ | ✅ |
| hover | ✅ | ✅ | ✅ | ✅ | ✅ |
| rename_symbol | ✅ | ✅ | ✅ | ✅ | ✅ |
| code_actions | ✅ | ✅ | ✅ | ✅ | ⚠️ |
| code_lens | ✅ | ✅ | ⚠️ | ⚠️ | ❌ |
| signature_help | ✅ | ✅ | ✅ | ✅ | ✅ |
| document_symbols | ✅ | ✅ | ✅ | ✅ | ✅ |
| call_hierarchy | ✅ | ✅ | ✅ | ✅ | ⚠️ |

Legend:
- ✅ Fully supported
- ⚠️ Partial support / version-dependent
- ❌ Not supported

*Note: Matrix populated from manual testing. Your mileage may vary based on LSP server versions.*

### Known Limitations

1. **Dynamic Capability Updates**: Tool availability is determined at server startup based on the capabilities advertised in the `initialize` response. Dynamic capability registration (capabilities added after startup) is not currently supported.

2. **Graceful Degradation**: If a tool is registered but the LSP server fails to handle the request (buggy server), the tool will return an error to the MCP client rather than failing silently.

## Tools

- `definition`: Retrieves the complete source code definition of any symbol (function, type, constant, etc.) from your codebase.
- `references`: Locates all usages and references of a symbol throughout the codebase.
- `diagnostics`: Provides diagnostic information for a specific file, including warnings and errors.
- `hover`: Display documentation, type hints, or other hover information for a given location.
- `rename_symbol`: Rename a symbol across a project.
- `edit_file`: Allows making multiple text edits to a file based on line numbers. Provides a more reliable and context-economical way to edit files compared to search and replace based edit tools.

## About

This codebase makes use of edited code from [gopls](https://go.googlesource.com/tools/+/refs/heads/master/gopls/internal/protocol) to handle LSP communication. See ATTRIBUTION for details. Everything here is covered by a permissive BSD style license.

[mcp-go](https://github.com/mark3labs/mcp-go) is used for MCP communication. Thank you for your service.

This is beta software. Please let me know by creating an issue if you run into any problems or have suggestions of any kind.

## Contributing

Please keep PRs small and open Issues first for anything substantial. AI slop O.K. as long as it is tested, passes checks, and doesn't smell too bad.

### Setup

Clone the repo:

```bash
git clone https://github.com/isaacphi/mcp-language-server.git
cd mcp-language-server
```

A [justfile](https://just.systems/man/en/) is included for convenience:

```bash
just -l
Available recipes:
    build    # Build
    check    # Run code audit checks
    fmt      # Format code
    generate # Generate LSP types and methods
    help     # Help
    install  # Install locally
    snapshot # Update snapshot tests
    test     # Run tests
```

Configure your Claude Desktop (or similar) to use the local binary:

```json
{
  "mcpServers": {
    "language-server": {
      "command": "/full/path/to/your/clone/mcp-language-server/mcp-language-server",
      "args": [
        "--workspace",
        "/path/to/workspace",
        "--lsp",
        "language-server-executable"
      ],
      "env": {
        "LOG_LEVEL": "DEBUG"
      }
    }
  }
}
```

Rebuild after making changes.

### Logging

Setting the `LOG_LEVEL` environment variable to DEBUG enables verbose logging to stderr for all components including messages to and from the language server and the language server's logs.

### LSP interaction

- `internal/lsp/methods.go` contains generated code to make calls to the connected language server.
- `internal/protocol/tsprotocol.go` contains generated code for LSP types. I borrowed this from `gopls`'s source code. Thank you for your service.
- LSP allows language servers to return different types for the same methods. Go doesn't like this so there are some ugly workarounds in `internal/protocol/interfaces.go`.

### Local Development and Snapshot Tests

There is a snapshot test suite that makes it a lot easier to try out changes to tools. These run actual language servers on mock workspaces and capture output and logs.

You will need the language servers installed locally to run them. There are tests for go, rust, python, and typescript.

```
integrationtests/
├── tests/        # Tests are in this folder
├── snapshots/    # Snapshots of tool outputs
├── test-output/  # Gitignored folder showing the final state of each workspace and logs after each test run
└── workspaces/   # Mock workspaces that the tools run on
```

To update snapshots, run `UPDATE_SNAPSHOTS=true go test ./integrationtests/...`
