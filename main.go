package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/isaacphi/mcp-language-server/internal/fileops"
	"github.com/isaacphi/mcp-language-server/internal/logging"
	"github.com/isaacphi/mcp-language-server/internal/lsp"
	"github.com/isaacphi/mcp-language-server/internal/protocol"
	"github.com/isaacphi/mcp-language-server/internal/watcher"
	"github.com/mark3labs/mcp-go/server"
)

// Create a logger for the core component
var coreLogger = logging.NewLogger(logging.Core)

type config struct {
	workspaceDir string
	lspCommand   string
	lspArgs      []string
	transport    string // "stdio" or "http"
	httpPort     int    // Port for HTTP transport (default: 8080)
}

type mcpServer struct {
	config           config
	lspClient        *lsp.Client
	mcpServer        *server.MCPServer
	ctx              context.Context
	cancelFunc       context.CancelFunc
	workspaceWatcher *watcher.WorkspaceWatcher
	capabilities     *protocol.ServerCapabilities
	fileOpsHandler   *fileops.FileOperationsHandler
}

// mcpNotificationListener implements fileops.FileOperationsListener
// and sends MCP notifications for file events
type mcpNotificationListener struct {
	server *server.MCPServer
}

// OnFileEvent implements fileops.FileOperationsListener
func (l *mcpNotificationListener) OnFileEvent(event fileops.FileEvent) {
	if l.server == nil {
		coreLogger.Debug("MCP server not initialized, skipping notification")
		return
	}

	if len(event.URIs) == 0 {
		coreLogger.Debug("No URIs in event, skipping notification")
		return
	}

	params := map[string]interface{}{
		"type": event.Type.String(),
		"uris": event.URIs,
	}

	coreLogger.Info("Sending MCP notification: type=%s, uris=%d", event.Type.String(), len(event.URIs))

	// Send notification to all connected MCP clients
	l.server.SendNotificationToAllClients(
		"notifications/resources/updated",
		params,
	)
}

func parseConfig() (*config, error) {
	cfg := &config{}
	flag.StringVar(&cfg.workspaceDir, "workspace", "", "Path to workspace directory")
	flag.StringVar(&cfg.lspCommand, "lsp", "", "LSP command to run (args should be passed after --)")
	flag.StringVar(&cfg.transport, "transport", "stdio", "Transport type: stdio or http")
	flag.IntVar(&cfg.httpPort, "port", 8080, "Port for HTTP transport")
	flag.Parse()

	// Get remaining args after -- as LSP arguments
	cfg.lspArgs = flag.Args()

	// Validate transport
	if cfg.transport != "stdio" && cfg.transport != "http" {
		return nil, fmt.Errorf("invalid transport: %s (must be stdio or http)", cfg.transport)
	}

	// Validate port for HTTP mode
	if cfg.transport == "http" {
		if cfg.httpPort < 1 || cfg.httpPort > 65535 {
			return nil, fmt.Errorf("invalid port: %d (must be 1-65535)", cfg.httpPort)
		}
	}

	// Validate workspace directory
	if cfg.workspaceDir == "" {
		return nil, fmt.Errorf("workspace directory is required")
	}

	workspaceDir, err := filepath.Abs(cfg.workspaceDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for workspace: %v", err)
	}
	cfg.workspaceDir = workspaceDir

	if _, err := os.Stat(cfg.workspaceDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("workspace directory does not exist: %s", cfg.workspaceDir)
	}

	// Validate LSP command
	if cfg.lspCommand == "" {
		return nil, fmt.Errorf("LSP command is required")
	}

	if _, err := exec.LookPath(cfg.lspCommand); err != nil {
		return nil, fmt.Errorf("LSP command not found: %s", cfg.lspCommand)
	}

	return cfg, nil
}

func newServer(config *config) (*mcpServer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &mcpServer{
		config:     *config,
		ctx:        ctx,
		cancelFunc: cancel,
	}, nil
}

func (s *mcpServer) initializeLSP() error {
	if err := os.Chdir(s.config.workspaceDir); err != nil {
		return fmt.Errorf("failed to change to workspace directory: %v", err)
	}

	client, err := lsp.NewClient(s.config.lspCommand, s.config.lspArgs...)
	if err != nil {
		return fmt.Errorf("failed to create LSP client: %v", err)
	}
	s.lspClient = client
	s.workspaceWatcher = watcher.NewWorkspaceWatcher(client)

	initResult, err := client.InitializeLSPClient(s.ctx, s.config.workspaceDir)
	if err != nil {
		return fmt.Errorf("initialize failed: %v", err)
	}

	// Store capabilities for tool registration
	s.capabilities = &initResult.Capabilities

	coreLogger.Debug("Server capabilities: %+v", initResult.Capabilities)

	go s.workspaceWatcher.WatchWorkspace(s.ctx, s.config.workspaceDir)
	return client.WaitForServerReady(s.ctx)
}

func (s *mcpServer) start() error {
	if err := s.initializeLSP(); err != nil {
		return err
	}

	s.mcpServer = server.NewMCPServer(
		"MCP Language Server",
		"v0.0.2",
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Create and wire file operations handler
	s.fileOpsHandler = fileops.NewFileOperationsHandler()

	// Create MCP notification listener and register with handler
	mcpListener := &mcpNotificationListener{server: s.mcpServer}
	s.fileOpsHandler.RegisterListener(mcpListener)
	coreLogger.Info("MCP notification listener registered with FileOperationsHandler")

	// Connect file operations handler to LSP client
	if s.lspClient != nil {
		s.lspClient.SetFileOperationsHandler(s.fileOpsHandler)
		coreLogger.Info("FileOperationsHandler connected to LSP client")
	}

	err := s.registerTools(s.capabilities)
	if err != nil {
		return fmt.Errorf("tool registration failed: %v", err)
	}

	// Transport selection based on config
	switch s.config.transport {
	case "stdio":
		coreLogger.Info("Starting MCP server with stdio transport")
		return server.ServeStdio(s.mcpServer)

	case "http":
		coreLogger.Info("Starting MCP server with HTTP transport on port %d", s.config.httpPort)
		return s.serveHTTP()

	default:
		return fmt.Errorf("unsupported transport: %s", s.config.transport)
	}
}

func (s *mcpServer) serveHTTP() error {
	// Create StreamableHTTPServer with mcp-go v0.43.2 API
	httpServer := server.NewStreamableHTTPServer(
		s.mcpServer,
		server.WithEndpointPath("/mcp/v1"),
		server.WithStateful(true),
		server.WithHeartbeatInterval(30*time.Second),
	)

	// Bind to localhost only for security
	addr := fmt.Sprintf("localhost:%d", s.config.httpPort)
	coreLogger.Info("HTTP server listening on %s", addr)

	// Start blocks until error or shutdown
	err := httpServer.Start(addr)
	if err != nil {
		// Provide helpful error messages
		if strings.Contains(err.Error(), "address already in use") {
			return fmt.Errorf("port %d is already in use, try a different port with --port", s.config.httpPort)
		}
		if strings.Contains(err.Error(), "permission denied") {
			return fmt.Errorf("permission denied to bind port %d, try a port > 1024", s.config.httpPort)
		}
		return fmt.Errorf("HTTP server failed: %w", err)
	}

	return nil
}

func main() {
	coreLogger.Info("MCP Language Server starting")

	done := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	config, err := parseConfig()
	if err != nil {
		coreLogger.Fatal("%v", err)
	}

	server, err := newServer(config)
	if err != nil {
		coreLogger.Fatal("%v", err)
	}

	// Parent process monitoring channel
	parentDeath := make(chan struct{})

	// Monitor parent process termination
	// Claude desktop does not properly kill child processes for MCP servers
	go func() {
		ppid := os.Getppid()
		coreLogger.Debug("Monitoring parent process: %d", ppid)

		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				currentPpid := os.Getppid()
				if currentPpid != ppid && (currentPpid == 1 || ppid == 1) {
					coreLogger.Info("Parent process %d terminated (current ppid: %d), initiating shutdown", ppid, currentPpid)
					close(parentDeath)
					return
				}
			case <-done:
				return
			}
		}
	}()

	// Handle shutdown triggers
	go func() {
		select {
		case sig := <-sigChan:
			coreLogger.Info("Received signal %v in PID: %d", sig, os.Getpid())
			cleanup(server, done)
		case <-parentDeath:
			coreLogger.Info("Parent death detected, initiating shutdown")
			cleanup(server, done)
		}
	}()

	if err := server.start(); err != nil {
		coreLogger.Error("Server error: %v", err)
		cleanup(server, done)
		os.Exit(1)
	}

	<-done
	coreLogger.Info("Server shutdown complete for PID: %d", os.Getpid())
	os.Exit(0)
}

func cleanup(s *mcpServer, done chan struct{}) {
	coreLogger.Info("Cleanup initiated for PID: %d", os.Getpid())

	// Create a context with timeout for shutdown operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if s.lspClient != nil {
		coreLogger.Info("Closing open files")
		s.lspClient.CloseAllFiles(ctx)

		// Create a shorter timeout context for the shutdown request
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 500*time.Millisecond)
		defer shutdownCancel()

		// Run shutdown in a goroutine with timeout to avoid blocking if LSP doesn't respond
		shutdownDone := make(chan struct{})
		go func() {
			coreLogger.Info("Sending shutdown request")
			if err := s.lspClient.Shutdown(shutdownCtx); err != nil {
				coreLogger.Error("Shutdown request failed: %v", err)
			}
			close(shutdownDone)
		}()

		// Wait for shutdown with timeout
		select {
		case <-shutdownDone:
			coreLogger.Info("Shutdown request completed")
		case <-time.After(1 * time.Second):
			coreLogger.Warn("Shutdown request timed out, proceeding with exit")
		}

		coreLogger.Info("Sending exit notification")
		if err := s.lspClient.Exit(ctx); err != nil {
			coreLogger.Error("Exit notification failed: %v", err)
		}

		coreLogger.Info("Closing LSP client")
		if err := s.lspClient.Close(); err != nil {
			coreLogger.Error("Failed to close LSP client: %v", err)
		}
	}

	// Send signal to the done channel
	select {
	case <-done: // Channel already closed
	default:
		close(done)
	}

	coreLogger.Info("Cleanup completed for PID: %d", os.Getpid())
}
