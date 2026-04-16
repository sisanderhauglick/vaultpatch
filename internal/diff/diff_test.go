package diff

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/vault"
)

func TestDiff_Added(t *testing.T) {
	src := vault.SecretMap{}
	dst := vault.SecretMap{"NEW_KEY": "value"}
	changes := Diff(src, dst)
	if len(changes) != 1 || changes[0].Type != Added {
		t.Fatalf("expected 1 Added change, got %+v", changes)
	}
}

func TestDiff_Removed(t *testing.T) {
	src := vault.SecretMap{"OLD_KEY": "value"}
	dst := vault.SecretMap{}
	changes := Diff(src, dst)
	if len(changes) != 1 || changes[0].Type != Removed {
		t.Fatalf("expected 1 Removed change, got %+v", changes)
	}
}

func TestDiff_Modified(t *testing.T) {
	src := vault.SecretMap{"KEY": "old"}
	dst := vault.SecretMap{"KEY": "new"}
	changes := Diff(src, dst)
	if len(changes) != 1 || changes[0].Type != Modified {
		t.Fatalf("expected 1 Modified change, got %+v", changes)
	}
	if changes[0].OldVal != "old" || changes[0].NewVal != "new" {
		t.Errorf("unexpected values: %+v", changes[0])
	}
}

func TestDiff_NoChanges(t *testing.T) {
	sm := vault.SecretMap{"KEY": "value"}
	changes := Diff(sm, sm)
	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %+v", changes)
	}
}

func TestDiff_Mixed(t *testing.T) {
	src := vault.SecretMap{"A": "1", "B": "2"}
	dst := vault.SecretMap{"A": "99", "C": "3"}
	changes := Diff(src, dst)
	if len(changes) != 3 {
		t.Fatalf("expected 3 changes, got %d: %+v", len(changes), changes)
	}
}
