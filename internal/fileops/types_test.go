package fileops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileEventType_String(t *testing.T) {
	tests := []struct {
		name     string
		eventType FileEventType
		expected string
	}{
		{"Created", FileEventCreated, "created"},
		{"Renamed", FileEventRenamed, "renamed"},
		{"Deleted", FileEventDeleted, "deleted"},
		{"Unknown", FileEventType(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.eventType.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFileEvent_Structure(t *testing.T) {
	t.Run("CreateEvent", func(t *testing.T) {
		event := FileEvent{
			Type: FileEventCreated,
			URIs: []string{"file:///workspace/test.go"},
		}
		assert.Equal(t, FileEventCreated, event.Type)
		assert.Len(t, event.URIs, 1)
		assert.Equal(t, "file:///workspace/test.go", event.URIs[0])
	})

	t.Run("RenameEvent", func(t *testing.T) {
		event := FileEvent{
			Type: FileEventRenamed,
			URIs: []string{"file:///workspace/old.go", "file:///workspace/new.go"},
		}
		assert.Equal(t, FileEventRenamed, event.Type)
		assert.Len(t, event.URIs, 2)
		assert.Equal(t, "file:///workspace/old.go", event.URIs[0])
		assert.Equal(t, "file:///workspace/new.go", event.URIs[1])
	})

	t.Run("DeleteEvent", func(t *testing.T) {
		event := FileEvent{
			Type: FileEventDeleted,
			URIs: []string{"file:///workspace/deleted.go"},
		}
		assert.Equal(t, FileEventDeleted, event.Type)
		assert.Len(t, event.URIs, 1)
		assert.Equal(t, "file:///workspace/deleted.go", event.URIs[0])
	})
}

func TestFileCreate(t *testing.T) {
	fc := FileCreate{URI: "file:///workspace/test.go"}
	assert.Equal(t, "file:///workspace/test.go", fc.URI)
}

func TestFileRename(t *testing.T) {
	fr := FileRename{
		OldURI: "file:///workspace/old.go",
		NewURI: "file:///workspace/new.go",
	}
	assert.Equal(t, "file:///workspace/old.go", fr.OldURI)
	assert.Equal(t, "file:///workspace/new.go", fr.NewURI)
}

func TestFileDelete(t *testing.T) {
	fd := FileDelete{URI: "file:///workspace/deleted.go"}
	assert.Equal(t, "file:///workspace/deleted.go", fd.URI)
}
