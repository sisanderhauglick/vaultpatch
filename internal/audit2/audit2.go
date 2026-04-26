package audit2

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Entry represents a single structured audit log entry.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Operation string            `json:"operation"`
	Path      string            `json:"path"`
	DryRun    bool              `json:"dry_run"`
	Changes   map[string]string `json:"changes,omitempty"`
	Error     string            `json:"error,omitempty"`
}

// Logger writes structured audit entries to an io.Writer.
type Logger struct {
	w io.Writer
}

// NewLogger returns a Logger that writes to w.
func NewLogger(w io.Writer) (*Logger, error) {
	if w == nil {
		return nil, fmt.Errorf("audit2: writer must not be nil")
	}
	return &Logger{w: w}, nil
}

// Log serialises entry as a JSON line and writes it to the logger's writer.
func (l *Logger) Log(entry Entry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	b, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit2: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}

// LogError is a convenience wrapper that records an operation failure.
func (l *Logger) LogError(op, path string, dryRun bool, err error) error {
	return l.Log(Entry{
		Operation: op,
		Path:      path,
		DryRun:    dryRun,
		Error:     err.Error(),
	})
}
