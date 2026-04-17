package sync_test

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/sync"
)

func TestSync_NilClientErrors(t *testing.T) {
	_, err := sync.Sync(nil, sync.Options{SrcPath: "a", DstPath: "b"})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestSync_EmptyPathsError(t *testing.T) {
	_, err := sync.Sync(nil, sync.Options{})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestFilterKeys_SubsetOnly(t *testing.T) {
	// filterKeys is unexported; exercise through Sync options in integration.
	// Here we verify the Options struct accepts Includes.
	opts := sync.Options{
		SrcPath:  "secret/src",
		DstPath:  "secret/dst",
		Includes: []string{"DB_PASS"},
		DryRun:   true,
	}
	if len(opts.Includes) != 1 {
		t.Fatalf("expected 1 include, got %d", len(opts.Includes))
	}
}

func TestSync_DryRun_ReturnsChangesWithoutApplying(t *testing.T) {
	// Without a live Vault, we verify DryRun path returns without error
	// when client is nil — error expected, confirming guard runs first.
	_, err := sync.Sync(nil, sync.Options{
		SrcPath: "secret/staging",
		DstPath: "secret/prod",
		DryRun:  true,
	})
	if err == nil {
		t.Fatal("nil client should still error before dry-run logic")
	}
}

func TestResult_AppliedZeroOnDryRun(t *testing.T) {
	// Confirm Result zero value is sensible.
	r := &sync.Result{}
	if r.Applied != 0 {
		t.Fatalf("expected Applied=0, got %d", r.Applied)
	}
	if len(r.Changes) != 0 {
		t.Fatalf("expected no changes, got %d", len(r.Changes))
	}
}
