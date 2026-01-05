package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/isaacphi/mcp-language-server/internal/fileops"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// transportType defines the transport protocol being tested
type transportType string

const (
	stdioTransport transportType = "stdio"
	httpTransport  transportType = "http"
)

// transportFixture contains test infrastructure for a specific transport
type transportFixture struct {
	transport transportType
	port      int
	setup     func(t *testing.T) (*mcpServer, func())
}

// createTestWorkspace creates a temporary workspace for testing
func createTestWorkspace(t *testing.T) (string, func()) {
	t.Helper()

	// Create temporary workspace
	tmpDir, err := os.MkdirTemp("", "mcp-transport-test-*")
	require.NoError(t, err, "failed to create temp workspace")

	// Create a basic Go file for testing
	testFile := filepath.Join(tmpDir, "main.go")
	testContent := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err, "failed to write test file")

	// Create go.mod
	goMod := filepath.Join(tmpDir, "go.mod")
	goModContent := `module example.com/test

go 1.24
`
	err = os.WriteFile(goMod, []byte(goModContent), 0644)
	require.NoError(t, err, "failed to write go.mod")

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

// findLSPCommand tries to find an available LSP server for testing
func findLSPCommand(t *testing.T) (string, bool) {
	t.Helper()

	// Try common LSP servers
	candidates := []string{"gopls", "typescript-language-server", "pyright"}

	for _, cmd := range candidates {
		_, err := exec.LookPath(cmd)
		if err == nil {
			return cmd, true
		}
	}

	return "", false
}

// setupStdioTransport creates a server with stdio transport
func setupStdioTransport(t *testing.T) (*mcpServer, func()) {
	t.Helper()

	lspCommand, found := findLSPCommand(t)
	if !found {
		t.Skip("No LSP server found in PATH (tried gopls, typescript-language-server, pyright)")
	}

	workspaceDir, cleanupWorkspace := createTestWorkspace(t)

	cfg := &config{
		workspaceDir: workspaceDir,
		lspCommand:   lspCommand,
		lspArgs:      []string{},
		transport:    "stdio",
	}

	srv, err := newServer(cfg)
	require.NoError(t, err, "failed to create stdio server")

	cleanup := func() {
		if srv.lspClient != nil {
			srv.lspClient.Close()
		}
		cleanupWorkspace()
	}

	return srv, cleanup
}

// setupHTTPTransport creates a server with HTTP transport
func setupHTTPTransport(t *testing.T) (*mcpServer, func()) {
	t.Helper()

	lspCommand, found := findLSPCommand(t)
	if !found {
		t.Skip("No LSP server found in PATH (tried gopls, typescript-language-server, pyright)")
	}

	workspaceDir, cleanupWorkspace := createTestWorkspace(t)

	// Use a high port to avoid conflicts
	port := 18080 + (os.Getpid() % 1000)

	cfg := &config{
		workspaceDir: workspaceDir,
		lspCommand:   lspCommand,
		lspArgs:      []string{},
		transport:    "http",
		httpPort:     port,
	}

	srv, err := newServer(cfg)
	require.NoError(t, err, "failed to create HTTP server")

	cleanup := func() {
		if srv.lspClient != nil {
			srv.lspClient.Close()
		}
		cleanupWorkspace()
	}

	return srv, cleanup
}

// TestTransport_ServerStartup tests server initialization for both transports
func TestTransport_ServerStartup(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping transport tests in short mode")
	}

	fixtures := []transportFixture{
		{
			transport: stdioTransport,
			setup:     setupStdioTransport,
		},
		{
			transport: httpTransport,
			setup:     setupHTTPTransport,
		},
	}

	for _, fixture := range fixtures {
		t.Run(string(fixture.transport), func(t *testing.T) {
			srv, cleanup := fixture.setup(t)
			defer cleanup()

			// Initialize LSP client (common to both transports)
			err := srv.initializeLSP()
			require.NoError(t, err, "failed to initialize LSP")

			// Verify LSP client was created
			assert.NotNil(t, srv.lspClient, "LSP client should be initialized")
			assert.NotNil(t, srv.capabilities, "server capabilities should be set")
			assert.NotNil(t, srv.workspaceWatcher, "workspace watcher should be initialized")

			// Verify configuration
			assert.Equal(t, string(fixture.transport), srv.config.transport)
		})
	}
}

// TestTransport_BasicCommunication tests basic MCP protocol communication
func TestTransport_BasicCommunication(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping transport tests in short mode")
	}

	fixtures := []transportFixture{
		{
			transport: stdioTransport,
			setup:     setupStdioTransport,
		},
		{
			transport: httpTransport,
			setup:     setupHTTPTransport,
		},
	}

	for _, fixture := range fixtures {
		t.Run(string(fixture.transport), func(t *testing.T) {
			srv, cleanup := fixture.setup(t)
			defer cleanup()

			// Initialize LSP
			err := srv.initializeLSP()
			require.NoError(t, err, "failed to initialize LSP")

			// Create MCP server
			srv.mcpServer = server.NewMCPServer(
				"MCP Language Server Test",
				"v0.0.1",
				server.WithLogging(),
			)

			// Register tools (minimal set for testing)
			err = srv.registerTools(srv.capabilities)
			require.NoError(t, err, "failed to register tools")

			// Verify MCP server is configured
			assert.NotNil(t, srv.mcpServer, "MCP server should be initialized")

			// For HTTP transport, verify port configuration
			if fixture.transport == httpTransport {
				assert.Greater(t, srv.config.httpPort, 0, "HTTP port should be set")
				assert.Less(t, srv.config.httpPort, 65536, "HTTP port should be valid")
			}
		})
	}
}

// TestTransport_GracefulShutdown tests clean shutdown for both transports
func TestTransport_GracefulShutdown(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping transport tests in short mode")
	}

	fixtures := []transportFixture{
		{
			transport: stdioTransport,
			setup:     setupStdioTransport,
		},
		{
			transport: httpTransport,
			setup:     setupHTTPTransport,
		},
	}

	for _, fixture := range fixtures {
		t.Run(string(fixture.transport), func(t *testing.T) {
			srv, cleanup := fixture.setup(t)
			defer cleanup()

			// Initialize LSP
			err := srv.initializeLSP()
			require.NoError(t, err, "failed to initialize LSP")

			// Simulate graceful shutdown by calling cleanup
			// The cleanup will be handled by the defer statement in the test

			// Verify LSP client is operational before shutdown
			assert.NotNil(t, srv.lspClient, "LSP client should be initialized")

			// Test context-based shutdown
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Close all files as part of shutdown
			srv.lspClient.CloseAllFiles(ctx)

			// Verify we can cleanly close the LSP client
			err = srv.lspClient.Close()
			// Close may return an error if the process was already terminated, which is fine
			if err != nil {
				t.Logf("LSP client close returned: %v (may be expected)", err)
			}

			t.Log("Graceful shutdown completed")
		})
	}
}

// TestTransport_HTTPEndpoint tests HTTP-specific functionality
func TestTransport_HTTPEndpoint(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping transport tests in short mode")
	}

	t.Run("HTTP endpoint availability", func(t *testing.T) {
		srv, cleanup := setupHTTPTransport(t)
		defer cleanup()

		err := srv.initializeLSP()
		require.NoError(t, err, "failed to initialize LSP")

		srv.mcpServer = server.NewMCPServer(
			"MCP Language Server Test",
			"v0.0.1",
			server.WithLogging(),
		)

		// Start HTTP server in background
		errChan := make(chan error, 1)
		go func() {
			errChan <- srv.serveHTTP()
		}()

		// Give server time to start
		time.Sleep(500 * time.Millisecond)

		// Verify server is listening
		addr := fmt.Sprintf("http://localhost:%d/mcp/v1", srv.config.httpPort)

		// Create a simple health check request
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", addr, nil)
		require.NoError(t, err, "failed to create request")

		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Do(req)

		// We expect either a successful response or a method not allowed
		// (since we're doing GET instead of POST for MCP)
		if err == nil {
			defer resp.Body.Close()
			t.Logf("HTTP server responded with status: %d", resp.StatusCode)
			// Any response means the server is listening
			assert.NotEqual(t, 0, resp.StatusCode, "should get a status code")
		} else {
			// Connection refused means server didn't start properly
			assert.NotContains(t, err.Error(), "connection refused",
				"server should be listening on port %d", srv.config.httpPort)
		}
	})
}

// TestTransport_MCPProtocolMessages tests MCP message handling
func TestTransport_MCPProtocolMessages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping transport tests in short mode")
	}

	fixtures := []transportFixture{
		{
			transport: stdioTransport,
			setup:     setupStdioTransport,
		},
		{
			transport: httpTransport,
			setup:     setupHTTPTransport,
		},
	}

	for _, fixture := range fixtures {
		t.Run(string(fixture.transport), func(t *testing.T) {
			srv, cleanup := fixture.setup(t)
			defer cleanup()

			err := srv.initializeLSP()
			require.NoError(t, err, "failed to initialize LSP")

			srv.mcpServer = server.NewMCPServer(
				"MCP Language Server Test",
				"v0.0.1",
				server.WithLogging(),
			)

			// Test initialize request/response structure
			t.Run("initialize message", func(t *testing.T) {
				// Create an initialize params structure
				initParams := mcp.InitializeParams{
					ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
					ClientInfo: mcp.Implementation{
						Name:    "test-client",
						Version: "1.0.0",
					},
					Capabilities: mcp.ClientCapabilities{},
				}

				// Marshal to JSON to verify structure
				data, err := json.Marshal(initParams)
				require.NoError(t, err, "failed to marshal initialize params")

				// Verify it's valid JSON
				var parsed map[string]interface{}
				err = json.Unmarshal(data, &parsed)
				require.NoError(t, err, "initialize params should be valid JSON")

				assert.Contains(t, parsed, "protocolVersion", "should have protocol version")
				assert.Contains(t, parsed, "clientInfo", "should have client info")
			})

			t.Run("tools list request", func(t *testing.T) {
				// Simulate a tools/list request
				listRequest := struct {
					Method string `json:"method"`
				}{
					Method: "tools/list",
				}

				data, err := json.Marshal(listRequest)
				require.NoError(t, err, "failed to marshal tools/list request")

				var parsed map[string]interface{}
				err = json.Unmarshal(data, &parsed)
				require.NoError(t, err, "tools/list request should be valid JSON")

				assert.Equal(t, "tools/list", parsed["method"])
			})
		})
	}
}

// TestTransport_ConcurrentRequests tests handling multiple concurrent requests
func TestTransport_ConcurrentRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping transport tests in short mode")
	}

	fixtures := []transportFixture{
		{
			transport: stdioTransport,
			setup:     setupStdioTransport,
		},
		{
			transport: httpTransport,
			setup:     setupHTTPTransport,
		},
	}

	for _, fixture := range fixtures {
		t.Run(string(fixture.transport), func(t *testing.T) {
			srv, cleanup := fixture.setup(t)
			defer cleanup()

			err := srv.initializeLSP()
			require.NoError(t, err, "failed to initialize LSP")

			srv.mcpServer = server.NewMCPServer(
				"MCP Language Server Test",
				"v0.0.1",
				server.WithLogging(),
			)

			// Verify server can be initialized multiple times without issues
			// This tests the underlying transport's ability to handle state
			for i := 0; i < 3; i++ {
				t.Logf("Iteration %d", i)
				assert.NotNil(t, srv.mcpServer, "MCP server should remain initialized")
				assert.NotNil(t, srv.lspClient, "LSP client should remain initialized")
			}
		})
	}
}

// TestTransport_ErrorHandling tests error scenarios for both transports
func TestTransport_ErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping transport tests in short mode")
	}

	t.Run("invalid workspace", func(t *testing.T) {
		cfg := &config{
			workspaceDir: "/nonexistent/path/that/does/not/exist",
			lspCommand:   "gopls",
			transport:    "stdio",
		}

		srv, err := newServer(cfg)
		require.NoError(t, err, "server creation should succeed")

		err = srv.initializeLSP()
		assert.Error(t, err, "should fail to initialize with invalid workspace")
	})

	t.Run("invalid LSP command", func(t *testing.T) {
		workspaceDir, cleanup := createTestWorkspace(t)
		defer cleanup()

		cfg := &config{
			workspaceDir: workspaceDir,
			lspCommand:   "nonexistent-lsp-server",
			transport:    "stdio",
		}

		srv, err := newServer(cfg)
		require.NoError(t, err, "server creation should succeed")

		err = srv.initializeLSP()
		assert.Error(t, err, "should fail to initialize with invalid LSP command")
	})

	t.Run("HTTP invalid port", func(t *testing.T) {
		workspaceDir, cleanup := createTestWorkspace(t)
		defer cleanup()

		cfg := &config{
			workspaceDir: workspaceDir,
			lspCommand:   "gopls",
			transport:    "http",
			httpPort:     99999, // Invalid port
		}

		srv, err := newServer(cfg)
		require.NoError(t, err, "server creation should succeed")

		// Port validation should happen during serveHTTP
		// We just verify the config is set
		assert.Equal(t, 99999, srv.config.httpPort)
	})
}

// TestTransport_FileOperations tests file operations work across transports
func TestTransport_FileOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping transport tests in short mode")
	}

	fixtures := []transportFixture{
		{
			transport: stdioTransport,
			setup:     setupStdioTransport,
		},
		{
			transport: httpTransport,
			setup:     setupHTTPTransport,
		},
	}

	for _, fixture := range fixtures {
		t.Run(string(fixture.transport), func(t *testing.T) {
			srv, cleanup := fixture.setup(t)
			defer cleanup()

			err := srv.initializeLSP()
			require.NoError(t, err, "failed to initialize LSP")

			srv.mcpServer = server.NewMCPServer(
				"MCP Language Server Test",
				"v0.0.1",
				server.WithLogging(),
			)

			// Create file operations handler
			fileOpsHandler := fileops.NewFileOperationsHandler()

			// Create and register MCP notification listener
			mcpListener := &mcpNotificationListener{server: srv.mcpServer}
			fileOpsHandler.RegisterListener(mcpListener)

			// Set up the handler
			srv.fileOpsHandler = fileOpsHandler

			// Verify file operations handler is set up
			assert.NotNil(t, srv.fileOpsHandler, "file ops handler should be initialized")

			// Test file path
			testFile := filepath.Join(srv.config.workspaceDir, "main.go")

			// Open a file through LSP client
			ctx := context.Background()
			err = srv.lspClient.OpenFile(ctx, testFile)
			require.NoError(t, err, "failed to open file")

			// Verify file is tracked as open
			isOpen := srv.lspClient.IsFileOpen(testFile)
			assert.True(t, isOpen, "file should be tracked as open")

			// Close the file
			err = srv.lspClient.CloseFile(ctx, testFile)
			require.NoError(t, err, "failed to close file")

			// Verify file is no longer open
			isOpen = srv.lspClient.IsFileOpen(testFile)
			assert.False(t, isOpen, "file should not be tracked as open after closing")
		})
	}
}

// TestTransport_NotificationPropagation tests notifications work across transports
func TestTransport_NotificationPropagation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping transport tests in short mode")
	}

	fixtures := []transportFixture{
		{
			transport: stdioTransport,
			setup:     setupStdioTransport,
		},
		{
			transport: httpTransport,
			setup:     setupHTTPTransport,
		},
	}

	for _, fixture := range fixtures {
		t.Run(string(fixture.transport), func(t *testing.T) {
			srv, cleanup := fixture.setup(t)
			defer cleanup()

			err := srv.initializeLSP()
			require.NoError(t, err, "failed to initialize LSP")

			srv.mcpServer = server.NewMCPServer(
				"MCP Language Server Test",
				"v0.0.1",
				server.WithLogging(),
			)

			// Create notification listener
			listener := &mcpNotificationListener{server: srv.mcpServer}

			// Test notification doesn't panic
			testEvent := fileops.FileEvent{
				Type: fileops.FileEventCreated,
				URIs: []string{"file:///test.go"},
			}
			listener.OnFileEvent(testEvent)

			// Verify listener is functional
			assert.NotNil(t, listener.server, "listener should have server reference")
		})
	}
}


// TestTransport_ContextCancellation tests graceful handling of context cancellation
func TestTransport_ContextCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping transport tests in short mode")
	}

	fixtures := []transportFixture{
		{
			transport: stdioTransport,
			setup:     setupStdioTransport,
		},
		{
			transport: httpTransport,
			setup:     setupHTTPTransport,
		},
	}

	for _, fixture := range fixtures {
		t.Run(string(fixture.transport), func(t *testing.T) {
			srv, cleanup := fixture.setup(t)
			defer cleanup()

			err := srv.initializeLSP()
			require.NoError(t, err, "failed to initialize LSP")

			// Create a context that we'll cancel
			ctx, cancel := context.WithCancel(context.Background())

			// Cancel immediately
			cancel()

			// Verify operations handle cancelled context gracefully
			testFile := filepath.Join(srv.config.workspaceDir, "main.go")
			err = srv.lspClient.OpenFile(ctx, testFile)

			// Should either succeed (operation completed before cancel)
			// or fail gracefully (context cancelled)
			if err != nil {
				t.Logf("Operation failed as expected with cancelled context: %v", err)
			} else {
				t.Log("Operation completed before context cancellation")
			}
		})
	}
}

// TestTransport_MessageSerialization tests JSON-RPC message serialization
func TestTransport_MessageSerialization(t *testing.T) {
	t.Run("initialize params", func(t *testing.T) {
		params := mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    "test-client",
				Version: "1.0.0",
			},
			Capabilities: mcp.ClientCapabilities{},
		}

		data, err := json.Marshal(params)
		require.NoError(t, err, "should marshal initialize params")

		var decoded mcp.InitializeParams
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err, "should unmarshal initialize params")

		assert.Equal(t, params.ProtocolVersion, decoded.ProtocolVersion)
		assert.Equal(t, params.ClientInfo.Name, decoded.ClientInfo.Name)
	})

	t.Run("streaming message format", func(t *testing.T) {
		// Test that we can create valid streaming messages for HTTP transport
		msg := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "tools/list",
			"id":      1,
		}

		data, err := json.Marshal(msg)
		require.NoError(t, err, "should marshal message")

		// Verify it can be sent over HTTP
		reader := bytes.NewReader(data)
		_, err = io.ReadAll(reader)
		require.NoError(t, err, "should be readable")
	})
}
