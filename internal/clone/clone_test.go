package clone

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/vault"
)

func TestClone_NilClientErrors(t *testing.T) {
	_, err := Clone(nil, Options{Source: "a", Destination: "b"})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestClone_EmptySourceErrors(t *testing.T) {
	c, _ := vault.NewClient(vault.Params{Addr: "http://127.0.0.1:8200", Token: "t"})
	_, err := Clone(c, Options{Destination: "b"})
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestClone_EmptyDestinationErrors(t *testing.T) {
	c, _ := vault.NewClient(vault.Params{Addr: "http://127.0.0.1:8200", Token: "t"})
	_, err := Clone(c, Options{Source: "a"})
	if err == nil {
		t.Fatal("expected error for empty destination")
	}
}

func TestFilterKeys_AllKeys(t *testing.T) {
	m := map[string]string{"a": "1", "b": "2"}
	out := filterKeys(m, nil)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestFilterKeys_SubsetOnly(t *testing.T) {
	m := map[string]string{"a": "1", "b": "2", "c": "3"}
	out := filterKeys(m, []string{"a", "c"})
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["b"]; ok {
		t.Fatal("key b should not be present")
	}
}

func TestFilterKeys_MissingKeyIgnored(t *testing.T) {
	m := map[string]string{"a": "1"}
	out := filterKeys(m, []string{"a", "z"})
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
}
