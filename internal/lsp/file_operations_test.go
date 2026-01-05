package lsp

import (
	"encoding/json"
	"testing"

	"github.com/isaacphi/mcp-language-server/internal/fileops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockFileOpsHandler implements a simple mock for testing
type mockFileOpsHandler struct {
	creates [][]fileops.FileCreate
	renames [][]fileops.FileRename
	deletes [][]fileops.FileDelete
}

func (m *mockFileOpsHandler) OnCreate(files []fileops.FileCreate) {
	m.creates = append(m.creates, files)
}

func (m *mockFileOpsHandler) OnRename(files []fileops.FileRename) {
	m.renames = append(m.renames, files)
}

func (m *mockFileOpsHandler) OnDelete(files []fileops.FileDelete) {
	m.deletes = append(m.deletes, files)
}

func TestClient_SetFileOperationsHandler(t *testing.T) {
	client := &Client{}
	handler := fileops.NewFileOperationsHandler()

	client.SetFileOperationsHandler(handler)

	assert.NotNil(t, client.fileOpsHandler)
}

func TestClient_HandleDidCreateFiles(t *testing.T) {
	client := &Client{}
	mock := &mockFileOpsHandler{}
	client.fileOpsHandler = mock

	params := map[string]interface{}{
		"files": []map[string]string{
			{"uri": "file:///workspace/test.go"},
		},
	}

	data, err := json.Marshal(params)
	require.NoError(t, err)

	err = client.handleDidCreateFiles(data)
	require.NoError(t, err)

	require.Len(t, mock.creates, 1)
	assert.Len(t, mock.creates[0], 1)
	assert.Equal(t, "file:///workspace/test.go", mock.creates[0][0].URI)
}

func TestClient_HandleDidCreateFiles_InvalidJSON(t *testing.T) {
	client := &Client{}
	mock := &mockFileOpsHandler{}
	client.fileOpsHandler = mock

	err := client.handleDidCreateFiles([]byte("invalid json"))

	// Should not error, just log
	assert.NoError(t, err)
	assert.Len(t, mock.creates, 0)
}

func TestClient_HandleDidRenameFiles(t *testing.T) {
	client := &Client{}
	mock := &mockFileOpsHandler{}
	client.fileOpsHandler = mock

	params := map[string]interface{}{
		"files": []map[string]string{
			{
				"oldUri": "file:///workspace/old.go",
				"newUri": "file:///workspace/new.go",
			},
		},
	}

	data, err := json.Marshal(params)
	require.NoError(t, err)

	err = client.handleDidRenameFiles(data)
	require.NoError(t, err)

	require.Len(t, mock.renames, 1)
	assert.Len(t, mock.renames[0], 1)
	assert.Equal(t, "file:///workspace/old.go", mock.renames[0][0].OldURI)
	assert.Equal(t, "file:///workspace/new.go", mock.renames[0][0].NewURI)
}

func TestClient_HandleDidDeleteFiles(t *testing.T) {
	client := &Client{}
	mock := &mockFileOpsHandler{}
	client.fileOpsHandler = mock

	params := map[string]interface{}{
		"files": []map[string]string{
			{"uri": "file:///workspace/test.go"},
		},
	}

	data, err := json.Marshal(params)
	require.NoError(t, err)

	err = client.handleDidDeleteFiles(data)
	require.NoError(t, err)

	require.Len(t, mock.deletes, 1)
	assert.Len(t, mock.deletes[0], 1)
	assert.Equal(t, "file:///workspace/test.go", mock.deletes[0][0].URI)
}

func TestClient_HandleWillCreateFiles_ReturnsNil(t *testing.T) {
	client := &Client{}

	params := map[string]interface{}{
		"files": []map[string]string{
			{"uri": "file:///workspace/test.go"},
		},
	}

	data, err := json.Marshal(params)
	require.NoError(t, err)

	result, err := client.handleWillCreateFiles(data)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestClient_HandleWillRenameFiles_ReturnsNil(t *testing.T) {
	client := &Client{}

	params := map[string]interface{}{
		"files": []map[string]string{
			{
				"oldUri": "file:///workspace/old.go",
				"newUri": "file:///workspace/new.go",
			},
		},
	}

	data, err := json.Marshal(params)
	require.NoError(t, err)

	result, err := client.handleWillRenameFiles(data)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestClient_HandleWillDeleteFiles_ReturnsNil(t *testing.T) {
	client := &Client{}

	params := map[string]interface{}{
		"files": []map[string]string{
			{"uri": "file:///workspace/test.go"},
		},
	}

	data, err := json.Marshal(params)
	require.NoError(t, err)

	result, err := client.handleWillDeleteFiles(data)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestClient_HandleFileOperations_NilHandler(t *testing.T) {
	client := &Client{}
	// Don't set fileOpsHandler

	params := map[string]interface{}{
		"files": []map[string]string{
			{"uri": "file:///workspace/test.go"},
		},
	}

	data, err := json.Marshal(params)
	require.NoError(t, err)

	// Should not panic when handler is nil
	require.NotPanics(t, func() {
		client.handleDidCreateFiles(data)
		client.handleDidRenameFiles(data)
		client.handleDidDeleteFiles(data)
	})
}
