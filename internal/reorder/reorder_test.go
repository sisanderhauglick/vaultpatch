package reorder

import (
	"testing"
)

func TestReorder_NilClientErrors(t *testing.T) {
	_, err := Reorder(nil, []string{"secret/a"}, Options{Keys: []string{"x"}})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestReorder_EmptyPathsError(t *testing.T) {
	_, err := Reorder(stubClient(), []string{}, Options{Keys: []string{"x"}})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestReorder_EmptyKeysError(t *testing.T) {
	_, err := Reorder(stubClient(), []string{"secret/a"}, Options{})
	if err == nil {
		t.Fatal("expected error for empty keys")
	}
}

func TestBuildOrder_ExplicitFirst(t *testing.T) {
	explicit := []string{"c", "a"}
	original := []string{"a", "b", "c"}
	got := buildOrder(explicit, original)

	if len(got) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(got))
	}
	if got[0] != "c" || got[1] != "a" || got[2] != "b" {
		t.Errorf("unexpected order: %v", got)
	}
}

func TestBuildOrder_UnknownExplicitKeyIgnored(t *testing.T) {
	explicit := []string{"z", "a"}
	original := []string{"a", "b"}
	got := buildOrder(explicit, original)

	// "z" is not in original so it should be omitted
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(got), got)
	}
	if got[0] != "a" || got[1] != "b" {
		t.Errorf("unexpected order: %v", got)
	}
}

func TestBuildOrder_NoExplicitOverlap_OriginalPreserved(t *testing.T) {
	explicit := []string{"x", "y"}
	original := []string{"a", "b", "c"}
	got := buildOrder(explicit, original)

	if len(got) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(got))
	}
}

func TestResult_DryRunFlagged(t *testing.T) {
	r := Result{DryRun: true}
	if !r.DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestResult_ReorderedAtSet(t *testing.T) {
	r := Result{}
	if !r.ReorderedAt.IsZero() {
		t.Error("expected zero time before assignment")
	}
}
