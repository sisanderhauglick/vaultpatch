package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/vaultpatch/vaultpatch/internal/diff"
)

// Entry represents a single audit log entry for a secret change.
type Entry struct {
	Timestamp time.Time       `json:"timestamp"`
	Path      string          `json:"path"`
	Changes   []diff.Change   `json:"changes"`
	DryRun    bool            `json:"dry_run"`
	Applied   bool            `json:"applied"`
}

// Logger writes audit entries to an io.Writer.
type Logger struct {
	w io.Writer
}

// NewLogger creates a new audit Logger writing to w.
func NewLogger(w io.Writer) *Logger {
	return &Logger{w: w}
}

// Log writes an audit entry for the given path and changes.
func (l *Logger) Log(path string, changes []diff.Change, dryRun bool, applied bool) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Path:      path,
		Changes:   changes,
		DryRun:    dryRun,
		Applied:   applied,
	}
	b, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}
