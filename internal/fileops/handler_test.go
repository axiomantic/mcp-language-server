package fileops

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockListener implements FileOperationsListener for testing
type mockListener struct {
	mu     sync.Mutex
	events []FileEvent
}

func (m *mockListener) OnFileEvent(event FileEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, event)
}

func (m *mockListener) GetEvents() []FileEvent {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]FileEvent{}, m.events...)
}

func TestNewFileOperationsHandler(t *testing.T) {
	handler := NewFileOperationsHandler()
	require.NotNil(t, handler)
}

func TestFileOperationsHandler_OnCreate(t *testing.T) {
	handler := NewFileOperationsHandler()
	listener := &mockListener{}
	handler.RegisterListener(listener)

	files := []FileCreate{
		{URI: "file:///workspace/test.go"},
		{URI: "file:///workspace/test2.go"},
	}

	handler.OnCreate(files)

	events := listener.GetEvents()
	require.Len(t, events, 1)
	assert.Equal(t, FileEventCreated, events[0].Type)
	assert.Equal(t, []string{"file:///workspace/test.go", "file:///workspace/test2.go"}, events[0].URIs)
}

func TestFileOperationsHandler_OnRename(t *testing.T) {
	handler := NewFileOperationsHandler()
	listener := &mockListener{}
	handler.RegisterListener(listener)

	files := []FileRename{
		{
			OldURI: "file:///workspace/old.go",
			NewURI: "file:///workspace/new.go",
		},
	}

	handler.OnRename(files)

	events := listener.GetEvents()
	require.Len(t, events, 1)
	assert.Equal(t, FileEventRenamed, events[0].Type)
	assert.Equal(t, []string{"file:///workspace/old.go", "file:///workspace/new.go"}, events[0].URIs)
}

func TestFileOperationsHandler_OnDelete(t *testing.T) {
	handler := NewFileOperationsHandler()
	listener := &mockListener{}
	handler.RegisterListener(listener)

	files := []FileDelete{
		{URI: "file:///workspace/test.go"},
	}

	handler.OnDelete(files)

	events := listener.GetEvents()
	require.Len(t, events, 1)
	assert.Equal(t, FileEventDeleted, events[0].Type)
	assert.Equal(t, []string{"file:///workspace/test.go"}, events[0].URIs)
}

func TestFileOperationsHandler_MultipleListeners(t *testing.T) {
	handler := NewFileOperationsHandler()
	listener1 := &mockListener{}
	listener2 := &mockListener{}

	handler.RegisterListener(listener1)
	handler.RegisterListener(listener2)

	files := []FileCreate{
		{URI: "file:///workspace/test.go"},
	}

	handler.OnCreate(files)

	// Both listeners should receive the event
	events1 := listener1.GetEvents()
	events2 := listener2.GetEvents()

	require.Len(t, events1, 1)
	require.Len(t, events2, 1)
	assert.Equal(t, events1[0], events2[0])
}

func TestFileOperationsHandler_UnregisterListener(t *testing.T) {
	handler := NewFileOperationsHandler()
	listener := &mockListener{}

	handler.RegisterListener(listener)
	handler.UnregisterListener(listener)

	files := []FileCreate{
		{URI: "file:///workspace/test.go"},
	}

	handler.OnCreate(files)

	// Listener should not receive event after unregistering
	events := listener.GetEvents()
	assert.Len(t, events, 0)
}

func TestFileOperationsHandler_NoListeners(t *testing.T) {
	handler := NewFileOperationsHandler()

	files := []FileCreate{
		{URI: "file:///workspace/test.go"},
	}

	// Should not panic when no listeners
	require.NotPanics(t, func() {
		handler.OnCreate(files)
	})
}

func TestFileOperationsHandler_EmptyFiles(t *testing.T) {
	handler := NewFileOperationsHandler()
	listener := &mockListener{}
	handler.RegisterListener(listener)

	// Test with empty slices
	handler.OnCreate([]FileCreate{})
	handler.OnRename([]FileRename{})
	handler.OnDelete([]FileDelete{})

	// Should not notify for empty file lists
	events := listener.GetEvents()
	assert.Len(t, events, 0)
}

func TestFileOperationsHandler_MultipleRenames(t *testing.T) {
	handler := NewFileOperationsHandler()
	listener := &mockListener{}
	handler.RegisterListener(listener)

	files := []FileRename{
		{OldURI: "file:///workspace/a.go", NewURI: "file:///workspace/a-renamed.go"},
		{OldURI: "file:///workspace/b.go", NewURI: "file:///workspace/b-renamed.go"},
	}

	handler.OnRename(files)

	events := listener.GetEvents()
	require.Len(t, events, 1)
	assert.Equal(t, FileEventRenamed, events[0].Type)
	// All old/new pairs should be in the URIs list
	expected := []string{
		"file:///workspace/a.go",
		"file:///workspace/a-renamed.go",
		"file:///workspace/b.go",
		"file:///workspace/b-renamed.go",
	}
	assert.Equal(t, expected, events[0].URIs)
}
