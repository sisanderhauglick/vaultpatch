package copy

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/vault"
)

func TestCopy_NilClientErrors(t *testing.T) {
	_, err := Copy(nil, Options{SourcePath: "a", DestPath: "b"})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestCopy_EmptySourceErrors(t *testing.T) {
	client := &vault.Client{}
	_, err := Copy(client, Options{SourcePath: "", DestPath: "b"})
	if err == nil {
		t.Fatal("expected error for empty source path")
	}
}

func TestCopy_EmptyDestErrors(t *testing.T) {
	client := &vault.Client{}
	_, err := Copy(client, Options{SourcePath: "a", DestPath: ""})
	if err == nil {
		t.Fatal("expected error for empty dest path")
	}
}

func TestFilterKeys_AllKeys(t *testing.T) {
	sm := vault.SecretMap{"a": "1", "b": "2"}
	out := filterKeys(sm, nil)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestFilterKeys_SubsetOnly(t *testing.T) {
	sm := vault.SecretMap{"a": "1", "b": "2", "c": "3"}
	out := filterKeys(sm, []string{"a", "c"})
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["b"]; ok {
		t.Fatal("key 'b' should have been filtered out")
	}
}

func TestFilterKeys_MissingKeySkipped(t *testing.T) {
	sm := vault.SecretMap{"a": "1"}
	out := filterKeys(sm, []string{"a", "z"})
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
}
