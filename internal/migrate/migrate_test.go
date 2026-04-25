package migrate

import (
	"testing"

	"github.com/yourusername/vaultpatch/internal/vault"
)

func TestMigrate_NilClientErrors(t *testing.T) {
	_, err := Migrate(nil, Options{Sources: []string{"a"}, Destination: "b"})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestMigrate_EmptySourcesErrors(t *testing.T) {
	c := &vault.Client{}
	_, err := Migrate(c, Options{Sources: []string{}, Destination: "b"})
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestMigrate_EmptyDestinationErrors(t *testing.T) {
	c := &vault.Client{}
	_, err := Migrate(c, Options{Sources: []string{"a"}, Destination: ""})
	if err == nil {
		t.Fatal("expected error for empty destination")
	}
}

func TestRemapKeys_IdentityWhenNoMap(t *testing.T) {
	src := map[string]string{"foo": "bar", "baz": "qux"}
	out := remapKeys(src, nil)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if out["foo"] != "bar" {
		t.Errorf("expected foo=bar, got %q", out["foo"])
	}
}

func TestRemapKeys_RenamesMatchedKeys(t *testing.T) {
	src := map[string]string{"old_key": "value", "keep": "same"}
	keyMap := map[string]string{"old_key": "new_key"}
	out := remapKeys(src, keyMap)
	if _, ok := out["old_key"]; ok {
		t.Error("old_key should have been renamed")
	}
	if out["new_key"] != "value" {
		t.Errorf("expected new_key=value, got %q", out["new_key"])
	}
	if out["keep"] != "same" {
		t.Errorf("expected keep=same, got %q", out["keep"])
	}
}

func TestResult_DryRunFlagged(t *testing.T) {
	r := Result{DryRun: true}
	if !r.DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestResult_MigratedAtSet(t *testing.T) {
	r := Result{}
	if !r.MigratedAt.IsZero() == false {
		t.Error("MigratedAt should be zero when unset")
	}
}
