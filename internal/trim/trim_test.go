package trim_test

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/trim"
)

func TestTrim_NilClientErrors(t *testing.T) {
	_, err := trim.Trim(nil, trim.Options{Paths: []string{"secret/a"}})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestTrim_EmptyPathsError(t *testing.T) {
	_, err := trim.Trim(stubClient(), trim.Options{})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestTrim_DryRun_ReturnsRemovedWithoutWrite(t *testing.T) {
	c := stubClient()
	results, err := trim.Trim(c, trim.Options{
		Paths:  []string{"secret/env"},
		Keys:   []string{"OLD_KEY"},
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].DryRun {
		t.Error("expected DryRun flag set")
	}
	if len(results[0].Removed) != 1 || results[0].Removed[0] != "OLD_KEY" {
		t.Errorf("unexpected removed keys: %v", results[0].Removed)
	}
}

func TestTrim_AllKeysWhenNoneSpecified(t *testing.T) {
	c := stubClient()
	results, err := trim.Trim(c, trim.Options{
		Paths:  []string{"secret/env"},
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results[0].Removed) == 0 {
		t.Error("expected all keys to be selected")
	}
}

func TestTrim_NoMatchingKeys_EmptyResult(t *testing.T) {
	c := stubClient()
	results, err := trim.Trim(c, trim.Options{
		Paths:  []string{"secret/env"},
		Keys:   []string{"NONEXISTENT"},
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

// stubClient returns a minimal vault.Client substitute via a local fake.
func stubClient() *vaultStub { return &vaultStub{} }

// vaultStub satisfies the interface expected by trim (duck-typed in tests).
type vaultStub struct{}
