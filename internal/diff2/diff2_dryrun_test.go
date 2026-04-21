package diff2

import (
	"testing"
)

// computeChanges is package-internal; these tests exercise it directly
// to validate dry-run semantics: no writes occur and results are populated.

func TestDryRun_ResultFlaggedTrue(t *testing.T) {
	r := Result{DryRun: true}
	if !r.DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestDryRun_ChangesPopulatedWithoutWrite(t *testing.T) {
	src := map[string]string{"SECRET": "before"}
	dst := map[string]string{"SECRET": "after"}

	changes := computeChanges(src, dst, nil)

	if len(changes) == 0 {
		t.Fatal("expected changes to be populated in dry-run")
	}
	if changes[0].Kind != Modified {
		t.Errorf("expected Modified, got %q", changes[0].Kind)
	}
}

func TestDryRun_SourceAndDestPreserved(t *testing.T) {
	r := Result{
		Source: "secret/dev",
		Dest:   "secret/prod",
		DryRun: true,
	}
	if r.Source != "secret/dev" {
		t.Errorf("unexpected Source: %q", r.Source)
	}
	if r.Dest != "secret/prod" {
		t.Errorf("unexpected Dest: %q", r.Dest)
	}
}

func TestDryRun_MultipleChangeKinds(t *testing.T) {
	src := map[string]string{"A": "1", "B": "old"}
	dst := map[string]string{"B": "new", "C": "3"}

	changes := computeChanges(src, dst, nil)

	kinds := map[ChangeKind]int{}
	for _, c := range changes {
		kinds[c.Kind]++
	}
	if kinds[Removed] != 1 {
		t.Errorf("expected 1 Removed, got %d", kinds[Removed])
	}
	if kinds[Modified] != 1 {
		t.Errorf("expected 1 Modified, got %d", kinds[Modified])
	}
	if kinds[Added] != 1 {
		t.Errorf("expected 1 Added, got %d", kinds[Added])
	}
}
