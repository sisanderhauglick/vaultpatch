package truncate_test

import (
	"testing"

	"github.com/youorg/vaultpatch/internal/truncate"
	"github.com/youorg/vaultpatch/internal/vault"
)

func TestTruncate_NilClientErrors(t *testing.T) {
	_, err := truncate.Truncate(nil, truncate.Options{Paths: []string{"secret/a"}, MaxLen: 10})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestTruncate_EmptyPathsError(t *testing.T) {
	c, _ := vault.NewClient(vault.Params{Addr: "http://127.0.0.1:8200", Token: "tok"})
	_, err := truncate.Truncate(c, truncate.Options{MaxLen: 10})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestTruncate_ZeroMaxLenError(t *testing.T) {
	c, _ := vault.NewClient(vault.Params{Addr: "http://127.0.0.1:8200", Token: "tok"})
	_, err := truncate.Truncate(c, truncate.Options{Paths: []string{"secret/a"}, MaxLen: 0})
	if err == nil {
		t.Fatal("expected error for zero MaxLen")
	}
}

func TestTruncate_DryRun_ResultFlagged(t *testing.T) {
	sc := newStubClient(map[string]interface{}{
		"token": "abcdefghij_extra",
	})
	results, err := truncate.Truncate(sc, truncate.Options{
		Paths:  []string{"secret/svc"},
		MaxLen: 10,
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].DryRun {
		t.Error("expected DryRun to be true")
	}
	if _, ok := results[0].Truncated["token"]; !ok {
		t.Error("expected 'token' to appear in Truncated map")
	}
}

func TestTruncate_ShortValue_Unchanged(t *testing.T) {
	sc := newStubClient(map[string]interface{}{
		"key": "short",
	})
	results, err := truncate.Truncate(sc, truncate.Options{
		Paths:  []string{"secret/svc"},
		MaxLen: 20,
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results[0].Truncated) != 0 {
		t.Errorf("expected no truncations, got %d", len(results[0].Truncated))
	}
}

func TestTruncate_SuffixAppended(t *testing.T) {
	sc := newStubClient(map[string]interface{}{
		"desc": "this is a very long description value",
	})
	results, err := truncate.Truncate(sc, truncate.Options{
		Paths:  []string{"secret/svc"},
		MaxLen: 7,
		Suffix: "...",
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := results[0].Truncated["desc"]
	expected := "this is..."
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestTruncate_ExactMaxLen_Unchanged(t *testing.T) {
	// A value whose length exactly equals MaxLen should not be truncated.
	sc := newStubClient(map[string]interface{}{
		"key": "exactly10c",
	})
	results, err := truncate.Truncate(sc, truncate.Options{
		Paths:  []string{"secret/svc"},
		MaxLen: 10,
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results[0].Truncated) != 0 {
		t.Errorf("expected no truncations for value at exact MaxLen, got %d", len(results[0].Truncated))
	}
}

// newStubClient returns a *vault.Client that satisfies the interface used by
// Truncate by pre-loading a fake path with the provided data.
func newStubClient(data map[string]interface{}) *vault.Client {
	c, _ := vault.NewClient(vault.Params{Addr: "http://127.0.0.1:8200", Token: "root"})
	_ = vault.WriteSecrets(c, "secret/svc", data)
	return c
}
