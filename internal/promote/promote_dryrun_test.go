package promote_test

import (
	"testing"

	"github.com/vaultpatch/internal/diff"
	"github.com/vaultpatch/internal/promote"
)

// stubResult simulates what Promote returns for dry-run scenarios.
func stubResult(changes []diff.Change, dryRun bool) *promote.Result {
	r := &promote.Result{Changes: changes}
	for range changes {
		if dryRun {
			r.Skipped++
		} else {
			r.Applied++
		}
	}
	return r
}

func TestDryRun_SkipsAllChanges(t *testing.T) {
	changes := []diff.Change{
		{Key: "FOO", Action: diff.Added, NewValue: "bar"},
		{Key: "BAZ", Action: diff.Modified, OldValue: "old", NewValue: "new"},
	}
	r := stubResult(changes, true)
	if r.Applied != 0 {
		t.Errorf("expected 0 applied, got %d", r.Applied)
	}
	if r.Skipped != 2 {
		t.Errorf("expected 2 skipped, got %d", r.Skipped)
	}
}

func TestLiveRun_AppliesAllChanges(t *testing.T) {
	changes := []diff.Change{
		{Key: "FOO", Action: diff.Added, NewValue: "bar"},
	}
	r := stubResult(changes, false)
	if r.Applied != 1 {
		t.Errorf("expected 1 applied, got %d", r.Applied)
	}
	if r.Skipped != 0 {
		t.Errorf("expected 0 skipped, got %d", r.Skipped)
	}
}
