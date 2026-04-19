package patch_test

import (
	"testing"

	"github.com/youorg/vaultpatch/internal/patch"
	"github.com/youorg/vaultpatch/internal/vault"
)

func stubClient() *vault.Client {
	c, _ := vault.NewClient(vault.Params{
		Addr:  "http://127.0.0.1:8200",
		Token: "test-token",
	})
	return c
}

func TestDryRun_PatchedAtSet(t *testing.T) {
	ops := []patch.Op{{Key: "env", Value: "staging"}}
	results, err := patch.Patch(stubClient(), []string{"secret/cfg"}, ops, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].PatchedAt.IsZero() {
		t.Error("expected PatchedAt to be set")
	}
}

func TestDryRun_MultiPath_AllResults(t *testing.T) {
	ops := []patch.Op{{Key: "x", Value: "1"}}
	paths := []string{"secret/a", "secret/b", "secret/c"}
	results, err := patch.Patch(stubClient(), paths, ops, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != len(paths) {
		t.Errorf("expected %d results, got %d", len(paths), len(results))
	}
}

func TestDryRun_PathPreservedInResult(t *testing.T) {
	ops := []patch.Op{{Key: "k", Value: "v"}}
	results, _ := patch.Patch(stubClient(), []string{"secret/mypath"}, ops, true)
	if results[0].Path != "secret/mypath" {
		t.Errorf("expected path %q, got %q", "secret/mypath", results[0].Path)
	}
}
