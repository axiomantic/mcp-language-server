package tools

import "time"

// MaxLineLength is the maximum line length for code formatting and display
const MaxLineLength = 9999

// DefaultContextLines is the default number of context lines to show around a match
const DefaultContextLines = 4

// DefaultNotificationTimeout is the default timeout for waiting for LSP notifications
const DefaultNotificationTimeout = 5 * time.Second

// FileOperationDebounceWindow is the time window for debouncing file operation notifications
const FileOperationDebounceWindow = 100 * time.Millisecond
