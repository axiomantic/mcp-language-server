package tools

import (
	"testing"
	"time"
)

func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		actual   interface{}
		expected interface{}
	}{
		{
			name:     "MaxLineLength should be 9999",
			actual:   MaxLineLength,
			expected: 9999,
		},
		{
			name:     "DefaultContextLines should be 4",
			actual:   DefaultContextLines,
			expected: 4,
		},
		{
			name:     "DefaultNotificationTimeout should be 5 seconds",
			actual:   DefaultNotificationTimeout,
			expected: 5 * time.Second,
		},
		{
			name:     "FileOperationDebounceWindow should be 100 milliseconds",
			actual:   FileOperationDebounceWindow,
			expected: 100 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("%s: got %v, want %v", tt.name, tt.actual, tt.expected)
			}
		})
	}
}
