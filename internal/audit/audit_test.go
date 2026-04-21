package audit_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/vaultpatch/vaultpatch/internal/audit"
	"github.com/vaultpatch/vaultpatch/internal/diff"
)

func TestLog_WritesJSONEntry(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf)

	changes := []diff.Change{
		{Key: "DB_PASS", Action: diff.ActionAdded, NewValue: "secret"},
	}

	if err := logger.Log("secret/app", changes, false, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to unmarshal audit entry: %v", err)
	}

	if entry.Path != "secret/app" {
		t.Errorf("expected path secret/app, got %s", entry.Path)
	}
	if len(entry.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(entry.Changes))
	}
	if entry.DryRun {
		t.Error("expected DryRun false")
	}
	if !entry.Applied {
		t.Error("expected Applied true")
	}
}

func TestLog_DryRun(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf)

	if err := logger.Log("secret/cfg", nil, true, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to unmarshal audit entry: %v", err)
	}

	if !entry.DryRun {
		t.Error("expected DryRun true")
	}
	if entry.Applied {
		t.Error("expected Applied false")
	}
}

func TestLog_EmptyChanges(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf)

	if err := logger.Log("secret/empty", []diff.Change{}, false, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestLog_TimestampPresent(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf)

	if err := logger.Log("secret/ts", nil, false, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to unmarshal audit entry: %v", err)
	}

	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp in audit entry")
	}
}
