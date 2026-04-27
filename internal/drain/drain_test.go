package drain_test

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/drain"
)

func TestDrain_NilClientErrors(t *testing.T) {
	_, err := drain.Drain(nil, drain.Options{Paths: []string{"secret/app"}})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestDrain_EmptyPathsError(t *testing.T) {
	c := newStubClient(map[string]map[string]interface{}{})
	_, err := drain.Drain(c, drain.Options{})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestDrain_DryRun_ResultFlagged(t *testing.T) {
	data := map[string]map[string]interface{}{
		"secret/app": {"DB_PASS": "s3cr3t", "API_KEY": "abc"},
	}
	c := newStubClient(data)
	results, err := drain.Drain(c, drain.Options{
		Paths:  []string{"secret/app"},
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].DryRun {
		t.Error("expected DryRun to be true")
	}
	// Data must be unchanged in dry-run
	if len(data["secret/app"]) != 2 {
		t.Error("dry-run must not modify source data")
	}
}

func TestDrain_DrainedKeysPopulated(t *testing.T) {
	c := newStubClient(map[string]map[string]interface{}{
		"secret/app": {"DB_PASS": "s3cr3t", "API_KEY": "abc", "KEEP": "yes"},
	})
	results, err := drain.Drain(c, drain.Options{
		Paths:    []string{"secret/app"},
		Preserve: []string{"KEEP"},
		DryRun:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results[0].Drained) != 2 {
		t.Errorf("expected 2 drained keys, got %d", len(results[0].Drained))
	}
	if len(results[0].Preserved) != 1 {
		t.Errorf("expected 1 preserved key, got %d", len(results[0].Preserved))
	}
}

func TestDrain_DrainedAtSet(t *testing.T) {
	c := newStubClient(map[string]map[string]interface{}{
		"secret/app": {"X": "1"},
	})
	results, err := drain.Drain(c, drain.Options{
		Paths:  []string{"secret/app"},
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].DrainedAt.IsZero() {
		t.Error("expected DrainedAt to be set")
	}
}

func TestDrain_PathPreservedInResult(t *testing.T) {
	c := newStubClient(map[string]map[string]interface{}{
		"secret/svc": {"TOKEN": "t"},
	})
	results, _ := drain.Drain(c, drain.Options{
		Paths:  []string{"secret/svc"},
		DryRun: true,
	})
	if results[0].Path != "secret/svc" {
		t.Errorf("expected path %q, got %q", "secret/svc", results[0].Path)
	}
}
