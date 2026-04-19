package patch_test

import (
	"testing"

	"github.com/youorg/vaultpatch/internal/patch"
)

func TestPatch_NilClientErrors(t *testing.T) {
	_, err := patch.Patch(nil, []string{"secret/a"}, []patch.Op{{Key: "k", Value: "v"}}, true)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestPatch_EmptyPathsError(t *testing.T) {
	_, err := patch.Patch(stubClient(), []string{}, []patch.Op{{Key: "k", Value: "v"}}, true)
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestPatch_EmptyOpsError(t *testing.T) {
	_, err := patch.Patch(stubClient(), []string{"secret/a"}, []patch.Op{}, true)
	if err == nil {
		t.Fatal("expected error for empty ops")
	}
}

func TestPatch_DryRun_ResultFlagged(t *testing.T) {
	ops := []patch.Op{{Key: "foo", Value: "bar"}}
	results, err := patch.Patch(stubClient(), []string{"secret/a"}, ops, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestPatch_AppliedOpsCount(t *testing.T) {
	ops := []patch.Op{
		{Key: "a", Value: "1"},
		{Key: "b", Value: "2"},
	}
	results, err := patch.Patch(stubClient(), []string{"secret/x"}, ops, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results[0].Applied) != 2 {
		t.Errorf("expected 2 applied ops, got %d", len(results[0].Applied))
	}
}

func TestPatch_DeleteOp_Flagged(t *testing.T) {
	ops := []patch.Op{{Key: "remove_me", Delete: true}}
	results, err := patch.Patch(stubClient(), []string{"secret/y"}, ops, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Applied[0].Delete {
		t.Error("expected delete op to be preserved in Applied")
	}
}
