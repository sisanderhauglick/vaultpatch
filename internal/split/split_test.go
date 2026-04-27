package split

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/vault"
)

func TestSplit_NilClientErrors(t *testing.T) {
	_, err := Split(nil, "secret/src", Options{
		Assignments: map[string][]string{"secret/dst": {"KEY"}},
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestSplit_EmptySourceErrors(t *testing.T) {
	c := &vault.Client{}
	_, err := Split(c, "", Options{
		Assignments: map[string][]string{"secret/dst": {"KEY"}},
	})
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestSplit_EmptyAssignmentsErrors(t *testing.T) {
	c := &vault.Client{}
	_, err := Split(c, "secret/src", Options{})
	if err == nil {
		t.Fatal("expected error for empty assignments")
	}
}

func TestResult_DryRunFlagged(t *testing.T) {
	r := Result{DryRun: true}
	if !r.DryRun {
		t.Fatal("expected DryRun to be true")
	}
}

func TestResult_SourceAndDestPreserved(t *testing.T) {
	r := Result{Source: "secret/src", Destination: "secret/dst"}
	if r.Source != "secret/src" {
		t.Errorf("unexpected source: %s", r.Source)
	}
	if r.Destination != "secret/dst" {
		t.Errorf("unexpected destination: %s", r.Destination)
	}
}

func TestResult_KeysPreserved(t *testing.T) {
	keys := []string{"A", "B", "C"}
	r := Result{Keys: keys}
	if len(r.Keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(r.Keys))
	}
}

func TestResult_SplitAtSet(t *testing.T) {
	r := Result{}
	if !r.SplitAt.IsZero() {
		t.Fatal("expected zero time on empty result")
	}
}
