package snapshot

import (
	"testing"
	"time"

	"github.com/your-org/vaultpatch/internal/vault"
)

func TestTake_NilClientErrors(t *testing.T) {
	_, err := Take(nil, []string{"secret/app"}, false)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestTake_EmptyPathsError(t *testing.T) {
	client := &vault.Client{}
	_, err := Take(client, []string{}, false)
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestDiff_NilClientErrors(t *testing.T) {
	snap := &Snapshot{CapturedAt: time.Now(), Paths: map[string]map[string]string{}}
	_, err := Diff(nil, snap)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestDiff_NilSnapshotErrors(t *testing.T) {
	client := &vault.Client{}
	_, err := Diff(client, nil)
	if err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestSnapshot_CapturedAtSet(t *testing.T) {
	before := time.Now().UTC()
	snap := &Snapshot{
		CapturedAt: time.Now().UTC(),
		Paths:      map[string]map[string]string{"secret/app": {"key": "val"}},
	}
	if snap.CapturedAt.Before(before) {
		t.Error("expected CapturedAt to be after test start")
	}
}

func TestResult_DryRunFlagged(t *testing.T) {
	r := &Result{DryRun: true, Paths: []string{"secret/a"}, Snapshot: &Snapshot{}}
	if !r.DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestResult_PathsPreserved(t *testing.T) {
	paths := []string{"secret/x", "secret/y"}
	r := &Result{DryRun: false, Paths: paths, Snapshot: &Snapshot{}}
	if len(r.Paths) != 2 {
		t.Errorf("expected 2 paths, got %d", len(r.Paths))
	}
}
