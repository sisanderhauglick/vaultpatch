package broadcast

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/vault"
)

func TestBroadcast_NilClientErrors(t *testing.T) {
	_, err := Broadcast(nil, "secret/src", []string{"secret/dst"}, Options{})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestBroadcast_EmptySourceErrors(t *testing.T) {
	c := &vault.Client{}
	_, err := Broadcast(c, "", []string{"secret/dst"}, Options{})
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestBroadcast_EmptyDestsErrors(t *testing.T) {
	c := &vault.Client{}
	_, err := Broadcast(c, "secret/src", []string{}, Options{})
	if err == nil {
		t.Fatal("expected error for empty destinations")
	}
}

func TestFilterKeys_AllKeysWhenNoneSpecified(t *testing.T) {
	data := map[string]string{"a": "1", "b": "2"}
	out := filterKeys(data, nil)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestFilterKeys_SubsetOnly(t *testing.T) {
	data := map[string]string{"a": "1", "b": "2", "c": "3"}
	out := filterKeys(data, []string{"a", "c"})
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["b"]; ok {
		t.Fatal("key 'b' should have been excluded")
	}
}

func TestFilterKeys_MissingKeySkipped(t *testing.T) {
	data := map[string]string{"a": "1"}
	out := filterKeys(data, []string{"a", "z"})
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
}
