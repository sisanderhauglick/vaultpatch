package diff2

import (
	"testing"
	"time"
)

func TestDiff2_NilClientErrors(t *testing.T) {
	_, err := Diff2(nil, Options{Source: "a", Dest: "b"})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestDiff2_EmptySourceErrors(t *testing.T) {
	_, err := Diff2(nil, Options{Dest: "b"})
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestDiff2_EmptyDestErrors(t *testing.T) {
	_, err := Diff2(nil, Options{Source: "a"})
	if err == nil {
		t.Fatal("expected error for empty dest")
	}
}

func TestComputeChanges_Added(t *testing.T) {
	src := map[string]string{}
	dst := map[string]string{"NEW_KEY": "v1"}
	changes := computeChanges(src, dst, nil)
	if len(changes) != 1 || changes[0].Kind != Added {
		t.Fatalf("expected 1 Added change, got %+v", changes)
	}
	if changes[0].NewValue != "v1" {
		t.Errorf("expected NewValue=v1, got %q", changes[0].NewValue)
	}
}

func TestComputeChanges_Removed(t *testing.T) {
	src := map[string]string{"OLD_KEY": "v1"}
	dst := map[string]string{}
	changes := computeChanges(src, dst, nil)
	if len(changes) != 1 || changes[0].Kind != Removed {
		t.Fatalf("expected 1 Removed change, got %+v", changes)
	}
	if changes[0].OldValue != "v1" {
		t.Errorf("expected OldValue=v1, got %q", changes[0].OldValue)
	}
}

func TestComputeChanges_Modified(t *testing.T) {
	src := map[string]string{"K": "old"}
	dst := map[string]string{"K": "new"}
	changes := computeChanges(src, dst, nil)
	if len(changes) != 1 || changes[0].Kind != Modified {
		t.Fatalf("expected 1 Modified change, got %+v", changes)
	}
	if changes[0].OldValue != "old" || changes[0].NewValue != "new" {
		t.Errorf("unexpected values: %+v", changes[0])
	}
}

func TestComputeChanges_NoChanges(t *testing.T) {
	src := map[string]string{"K": "v"}
	dst := map[string]string{"K": "v"}
	changes := computeChanges(src, dst, nil)
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(changes))
	}
}

func TestComputeChanges_KeyFilter(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"A": "9", "B": "9"}
	changes := computeChanges(src, dst, []string{"A"})
	if len(changes) != 1 || changes[0].Key != "A" {
		t.Fatalf("expected only key A in changes, got %+v", changes)
	}
}

func TestResult_HasChanges(t *testing.T) {
	r := Result{Changes: []Change{{Key: "X", Kind: Added}}}
	if !r.HasChanges() {
		t.Error("expected HasChanges to return true")
	}
}

func TestResult_DiffedAtSet(t *testing.T) {
	r := Result{DiffedAt: time.Now().UTC()}
	if r.DiffedAt.IsZero() {
		t.Error("expected DiffedAt to be set")
	}
}
