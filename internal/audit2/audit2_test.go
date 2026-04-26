package audit2_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpatch/internal/audit2"
)

func TestNewLogger_NilWriter(t *testing.T) {
	_, err := audit2.NewLogger(nil)
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestLog_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	l, _ := audit2.NewLogger(&buf)

	entry := audit2.Entry{
		Operation: "rotate",
		Path:      "secret/prod/db",
		DryRun:    false,
		Changes:   map[string]string{"password": "changed"},
	}
	if err := l.Log(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got audit2.Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if got.Operation != "rotate" {
		t.Errorf("operation: got %q, want %q", got.Operation, "rotate")
	}
	if got.Path != "secret/prod/db" {
		t.Errorf("path: got %q, want %q", got.Path, "secret/prod/db")
	}
}

func TestLog_TimestampAutoSet(t *testing.T) {
	var buf bytes.Buffer
	l, _ := audit2.NewLogger(&buf)

	before := time.Now().UTC()
	_ = l.Log(audit2.Entry{Operation: "sync", Path: "secret/x"})
	after := time.Now().UTC()

	var got audit2.Entry
	_ = json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got)

	if got.Timestamp.Before(before) || got.Timestamp.After(after) {
		t.Errorf("timestamp %v not within expected range [%v, %v]", got.Timestamp, before, after)
	}
}

func TestLog_DryRunFlagged(t *testing.T) {
	var buf bytes.Buffer
	l, _ := audit2.NewLogger(&buf)

	_ = l.Log(audit2.Entry{Operation: "apply", Path: "secret/y", DryRun: true})

	var got audit2.Entry
	_ = json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got)
	if !got.DryRun {
		t.Error("expected dry_run to be true")
	}
}

func TestLogError_SetsErrorField(t *testing.T) {
	var buf bytes.Buffer
	l, _ := audit2.NewLogger(&buf)

	_ = l.LogError("patch", "secret/z", false, errors.New("permission denied"))

	var got audit2.Entry
	_ = json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got)
	if !strings.Contains(got.Error, "permission denied") {
		t.Errorf("error field: got %q, want to contain %q", got.Error, "permission denied")
	}
}
