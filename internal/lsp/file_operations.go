package lsp

import (
	"encoding/json"

	"github.com/isaacphi/mcp-language-server/internal/fileops"
)

// FileOperationsHandler interface for LSP client file operations
type FileOperationsHandler interface {
	OnCreate(files []fileops.FileCreate)
	OnRename(files []fileops.FileRename)
	OnDelete(files []fileops.FileDelete)
}

// SetFileOperationsHandler sets the file operations handler for this client
func (c *Client) SetFileOperationsHandler(handler *fileops.FileOperationsHandler) {
	c.fileOpsHandler = handler
}

// LSP file operation parameter types
type createFilesParams struct {
	Files []fileops.FileCreate `json:"files"`
}

type renameFilesParams struct {
	Files []fileops.FileRename `json:"files"`
}

type deleteFilesParams struct {
	Files []fileops.FileDelete `json:"files"`
}

// handleWillCreateFiles handles workspace/willCreateFiles request
// Returns nil to allow operation without modification
func (c *Client) handleWillCreateFiles(params json.RawMessage) (interface{}, error) {
	var createParams createFilesParams
	if err := json.Unmarshal(params, &createParams); err != nil {
		lspLogger.Error("Failed to parse willCreateFiles params: %v", err)
		return nil, nil // Don't fail the LSP connection
	}

	// Currently we don't modify file operations
	// Return nil to allow the operation
	return nil, nil
}

// handleDidCreateFiles handles workspace/didCreateFiles notification
func (c *Client) handleDidCreateFiles(params json.RawMessage) error {
	var createParams createFilesParams
	if err := json.Unmarshal(params, &createParams); err != nil {
		lspLogger.Error("Failed to parse didCreateFiles params: %v", err)
		return nil // Don't fail the LSP connection
	}

	if c.fileOpsHandler != nil {
		c.fileOpsHandler.OnCreate(createParams.Files)
	}

	return nil
}

// handleWillRenameFiles handles workspace/willRenameFiles request
func (c *Client) handleWillRenameFiles(params json.RawMessage) (interface{}, error) {
	var renameParams renameFilesParams
	if err := json.Unmarshal(params, &renameParams); err != nil {
		lspLogger.Error("Failed to parse willRenameFiles params: %v", err)
		return nil, nil
	}

	return nil, nil
}

// handleDidRenameFiles handles workspace/didRenameFiles notification
func (c *Client) handleDidRenameFiles(params json.RawMessage) error {
	var renameParams renameFilesParams
	if err := json.Unmarshal(params, &renameParams); err != nil {
		lspLogger.Error("Failed to parse didRenameFiles params: %v", err)
		return nil
	}

	if c.fileOpsHandler != nil {
		c.fileOpsHandler.OnRename(renameParams.Files)
	}

	return nil
}

// handleWillDeleteFiles handles workspace/willDeleteFiles request
func (c *Client) handleWillDeleteFiles(params json.RawMessage) (interface{}, error) {
	var deleteParams deleteFilesParams
	if err := json.Unmarshal(params, &deleteParams); err != nil {
		lspLogger.Error("Failed to parse willDeleteFiles params: %v", err)
		return nil, nil
	}

	return nil, nil
}

// handleDidDeleteFiles handles workspace/didDeleteFiles notification
func (c *Client) handleDidDeleteFiles(params json.RawMessage) error {
	var deleteParams deleteFilesParams
	if err := json.Unmarshal(params, &deleteParams); err != nil {
		lspLogger.Error("Failed to parse didDeleteFiles params: %v", err)
		return nil
	}

	if c.fileOpsHandler != nil {
		c.fileOpsHandler.OnDelete(deleteParams.Files)
	}

	return nil
}
