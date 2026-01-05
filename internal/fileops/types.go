package fileops

// FileEventType represents the type of file operation
type FileEventType int

const (
	// FileEventCreated indicates a file was created
	FileEventCreated FileEventType = iota
	// FileEventRenamed indicates a file was renamed
	FileEventRenamed
	// FileEventDeleted indicates a file was deleted
	FileEventDeleted
)

// String returns the string representation of FileEventType
func (t FileEventType) String() string {
	switch t {
	case FileEventCreated:
		return "created"
	case FileEventRenamed:
		return "renamed"
	case FileEventDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}

// FileEvent represents a file operation event
type FileEvent struct {
	Type FileEventType
	URIs []string // For create/delete: single URI; for rename: [oldURI, newURI]
}

// FileCreate represents a file to be created (from LSP)
type FileCreate struct {
	URI string `json:"uri"`
}

// FileRename represents a file rename operation (from LSP)
type FileRename struct {
	OldURI string `json:"oldUri"`
	NewURI string `json:"newUri"`
}

// FileDelete represents a file to be deleted (from LSP)
type FileDelete struct {
	URI string `json:"uri"`
}

// FileOperationsListener defines the interface for objects that want to receive file operation events
type FileOperationsListener interface {
	OnFileEvent(event FileEvent)
}
