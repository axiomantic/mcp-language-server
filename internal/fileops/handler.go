package fileops

import (
	"sync"
)

// FileOperationsHandler manages file operation events and notifies registered listeners
type FileOperationsHandler struct {
	listeners []FileOperationsListener
	mu        sync.RWMutex
}

// NewFileOperationsHandler creates a new file operations handler
func NewFileOperationsHandler() *FileOperationsHandler {
	return &FileOperationsHandler{
		listeners: make([]FileOperationsListener, 0),
	}
}

// RegisterListener adds a listener to receive file operation events
func (h *FileOperationsHandler) RegisterListener(listener FileOperationsListener) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.listeners = append(h.listeners, listener)
}

// UnregisterListener removes a listener from receiving file operation events
func (h *FileOperationsHandler) UnregisterListener(listener FileOperationsListener) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Find and remove the listener
	for i, l := range h.listeners {
		// Use pointer comparison
		if l == listener {
			h.listeners = append(h.listeners[:i], h.listeners[i+1:]...)
			break
		}
	}
}

// OnCreate handles file creation events
func (h *FileOperationsHandler) OnCreate(files []FileCreate) {
	if len(files) == 0 {
		return
	}

	uris := make([]string, len(files))
	for i, file := range files {
		uris[i] = file.URI
	}

	event := FileEvent{
		Type: FileEventCreated,
		URIs: uris,
	}

	h.notifyListeners(event)
}

// OnRename handles file rename events
func (h *FileOperationsHandler) OnRename(files []FileRename) {
	if len(files) == 0 {
		return
	}

	// For renames, include both old and new URIs in pairs
	uris := make([]string, 0, len(files)*2)
	for _, file := range files {
		uris = append(uris, file.OldURI, file.NewURI)
	}

	event := FileEvent{
		Type: FileEventRenamed,
		URIs: uris,
	}

	h.notifyListeners(event)
}

// OnDelete handles file deletion events
func (h *FileOperationsHandler) OnDelete(files []FileDelete) {
	if len(files) == 0 {
		return
	}

	uris := make([]string, len(files))
	for i, file := range files {
		uris[i] = file.URI
	}

	event := FileEvent{
		Type: FileEventDeleted,
		URIs: uris,
	}

	h.notifyListeners(event)
}

// notifyListeners sends an event to all registered listeners
func (h *FileOperationsHandler) notifyListeners(event FileEvent) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, listener := range h.listeners {
		listener.OnFileEvent(event)
	}
}
