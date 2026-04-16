package apply_test

import (
	"testing"

	"github.com/yourusername/vaultpatch/internal/apply"
	"github.com/yourusername/vaultpatch/internal/diff"
)

func TestApply_DryRun_AllSucceed(t *testing.T) {
	changes := []diff.Entry{
		{Key: "FOO", Type: diff.Added, NewValue: "bar"},
		{Key: "BAZ", Type: diff.Modified, OldValue: "old", NewValue: "new"},
		{Key: "DEL", Type: diff.Removed, OldValue: "gone"},
	}

	// nil client is safe in dry-run because no network calls are made.
	results := apply.Apply(nil, "secret", "myapp", changes, apply.Options{DryRun: true})

	if len(results) != len(changes) {
		t.Fatalf("expected %d results, got %d", len(changes), len(results))
	}
	for _, r := range results {
		if !r.Success {
			t.Errorf("expected success for key %q, got err: %v", r.Key, r.Err)
		}
	}
}

func TestApply_DryRun_NoChanges(t *testing.T) {
	results := apply.Apply(nil, "secret", "myapp", []diff.Entry{}, apply.Options{DryRun: true})
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestApply_DryRun_ActionLabels(t *testing.T) {
	changes := []diff.Entry{
		{Key: "A", Type: diff.Added, NewValue: "1"},
		{Key: "B", Type: diff.Removed, OldValue: "2"},
	}

	results := apply.Apply(nil, "secret", "myapp", changes, apply.Options{DryRun: true})

	if results[0].Action != "added" {
		t.Errorf("expected action 'added', got %q", results[0].Action)
	}
	if results[1].Action != "removed" {
		t.Errorf("expected action 'removed', got %q", results[1].Action)
	}
}
