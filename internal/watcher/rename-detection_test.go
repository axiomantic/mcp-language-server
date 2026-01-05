package watcher

import (
	"sync"
	"testing"
	"time"

	"github.com/isaacphi/mcp-language-server/internal/fileops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockFileOpsHandler implements a simple handler for testing
type mockFileOpsHandler struct {
	mu      sync.Mutex
	creates [][]fileops.FileCreate
	renames [][]fileops.FileRename
	deletes [][]fileops.FileDelete
}

func (m *mockFileOpsHandler) OnCreate(files []fileops.FileCreate) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.creates = append(m.creates, files)
}

func (m *mockFileOpsHandler) OnRename(files []fileops.FileRename) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.renames = append(m.renames, files)
}

func (m *mockFileOpsHandler) OnDelete(files []fileops.FileDelete) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deletes = append(m.deletes, files)
}

func (m *mockFileOpsHandler) getCreates() [][]fileops.FileCreate {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.creates
}

func (m *mockFileOpsHandler) getRenames() [][]fileops.FileRename {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.renames
}

func (m *mockFileOpsHandler) getDeletes() [][]fileops.FileDelete {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.deletes
}

func TestWorkspaceWatcher_RenameDetection(t *testing.T) {
	// Setup
	watcher := &WorkspaceWatcher{
		pendingEvents: make(map[string]*pendingFileEvent),
	}

	mock := &mockFileOpsHandler{}

	// Simulate delete followed by create within 100ms
	oldURI := "file:///workspace/old.go"
	newURI := "file:///workspace/new.go"

	// Step 1: Delete event
	watcher.trackDeleteEvent(oldURI, mock)
	assert.Len(t, watcher.pendingEvents, 1)

	// Step 2: Create event within window (50ms)
	time.Sleep(50 * time.Millisecond)
	detected := watcher.checkForRename(newURI, mock)

	// Verify rename was detected
	assert.True(t, detected)

	// Wait a bit for async operations
	time.Sleep(10 * time.Millisecond)

	renames := mock.getRenames()
	assert.Len(t, renames, 1)
	assert.Equal(t, oldURI, renames[0][0].OldURI)
	assert.Equal(t, newURI, renames[0][0].NewURI)

	// Verify no separate delete
	deletes := mock.getDeletes()
	assert.Len(t, deletes, 0)

	// Verify pending events cleaned up
	assert.Len(t, watcher.pendingEvents, 0)
}

func TestWorkspaceWatcher_NoRenameDetection_Timeout(t *testing.T) {
	watcher := &WorkspaceWatcher{
		pendingEvents: make(map[string]*pendingFileEvent),
	}

	mock := &mockFileOpsHandler{}

	oldURI := "file:///workspace/old.go"
	newURI := "file:///workspace/new.go"

	// Delete event
	watcher.trackDeleteEvent(oldURI, mock)

	// Wait beyond debounce window
	time.Sleep(150 * time.Millisecond)

	// Process timeout
	watcher.processDeleteTimeout(oldURI, mock)

	// Should emit delete
	deletes := mock.getDeletes()
	assert.Len(t, deletes, 1)
	assert.Equal(t, oldURI, deletes[0][0].URI)

	// No rename
	renames := mock.getRenames()
	assert.Len(t, renames, 0)

	// Now create arrives - should be separate create
	detected := watcher.checkForRename(newURI, mock)
	assert.False(t, detected)
}

func TestWorkspaceWatcher_MultipleDeletes(t *testing.T) {
	watcher := &WorkspaceWatcher{
		pendingEvents: make(map[string]*pendingFileEvent),
	}

	mock := &mockFileOpsHandler{}

	// Track multiple deletes
	watcher.trackDeleteEvent("file:///workspace/a.go", mock)
	watcher.trackDeleteEvent("file:///workspace/b.go", mock)
	watcher.trackDeleteEvent("file:///workspace/c.go", mock)

	assert.Len(t, watcher.pendingEvents, 3)

	// Match one
	detected := watcher.checkForRename("file:///workspace/a-renamed.go", mock)
	assert.True(t, detected)
	assert.Len(t, watcher.pendingEvents, 2)

	// Others timeout
	time.Sleep(150 * time.Millisecond)
	watcher.processDeleteTimeout("file:///workspace/b.go", mock)
	watcher.processDeleteTimeout("file:///workspace/c.go", mock)

	deletes := mock.getDeletes()
	assert.Len(t, deletes, 2)
	assert.Len(t, watcher.pendingEvents, 0)
}

func TestWorkspaceWatcher_NoRenameDetection_RegularCreate(t *testing.T) {
	watcher := &WorkspaceWatcher{
		pendingEvents: make(map[string]*pendingFileEvent),
	}

	mock := &mockFileOpsHandler{}

	newURI := "file:///workspace/new.go"

	// Create without prior delete
	detected := watcher.checkForRename(newURI, mock)
	assert.False(t, detected)

	// No rename event
	renames := mock.getRenames()
	assert.Len(t, renames, 0)
}

func TestWorkspaceWatcher_RegularDelete(t *testing.T) {
	watcher := &WorkspaceWatcher{
		pendingEvents: make(map[string]*pendingFileEvent),
	}

	mock := &mockFileOpsHandler{}

	uri := "file:///workspace/deleted.go"

	// Delete event
	watcher.trackDeleteEvent(uri, mock)
	require.Len(t, watcher.pendingEvents, 1)

	// Wait for timeout
	time.Sleep(150 * time.Millisecond)
	watcher.processDeleteTimeout(uri, mock)

	// Should emit delete
	deletes := mock.getDeletes()
	assert.Len(t, deletes, 1)
	assert.Equal(t, uri, deletes[0][0].URI)

	// Verify cleanup
	assert.Len(t, watcher.pendingEvents, 0)
}
