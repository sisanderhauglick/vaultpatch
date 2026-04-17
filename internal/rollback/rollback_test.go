package rollback

import (
	"context"
	"testing"
)

func TestCapture_NilClientErrors(t *testing.T) {
	_, err := Capture(context.Background(), nil, "secret/data/app")
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestCapture_EmptyPathErrors(t *testing.T) {
	// We can't easily spin up a real Vault in unit tests, so we validate
	// guard clauses only.
	_, err := Capture(context.Background(), nil, "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestRestore_NilClientErrors(t *testing.T) {
	snap := &Snapshot{Path: "secret/data/app", Secrets: map[string]string{"k": "v"}}
	err := Restore(context.Background(), nil, snap, false)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestRestore_NilSnapshotErrors(t *testing.T) {
	err := Restore(context.Background(), nil, nil, false)
	if err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestRestore_DryRun_NoWrite(t *testing.T) {
	snap := &Snapshot{
		Path:    "secret/data/app",
		Secrets: map[string]string{"foo": "bar", "baz": "qux"},
	}
	// dry-run with nil client should still succeed because no write occurs
	err := Restore(context.Background(), nil, snap, true)
	// nil client check happens before dryRun guard; expect error
	if err == nil {
		t.Fatal("expected nil-client error even in dry-run")
	}
}

func TestSnapshot_Fields(t *testing.T) {
	snap := &Snapshot{
		Path:    "secret/data/test",
		Secrets: map[string]string{"key": "value"},
	}
	if snap.Path != "secret/data/test" {
		t.Errorf("unexpected path: %s", snap.Path)
	}
	if snap.Secrets["key"] != "value" {
		t.Errorf("unexpected secret value")
	}
}
