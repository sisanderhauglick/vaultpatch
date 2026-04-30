package broadcast

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/vault"
)

// stubClient returns a *vault.Client backed by an in-memory map so that
// ReadSecrets and WriteSecrets work without a live Vault instance.
type stubClient struct {
	store map[string]map[string]string
}

func newStub(data map[string]string) *vault.Client {
	// Re-use the package-level helper already present in vault tests.
	// In this test file we build a real *vault.Client via the exported
	// test helper so we can inject an in-memory backend.
	return vault.NewTestClient(data)
}

func TestDryRun_ResultFlagged(t *testing.T) {
	c := vault.NewTestClient(map[string]string{"key": "val"})
	results, err := Broadcast(c, "secret/src", []string{"secret/dst1", "secret/dst2"}, Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if !r.DryRun {
			t.Errorf("expected DryRun=true for result %+v", r)
		}
	}
}

func TestDryRun_ResultCountMatchesDests(t *testing.T) {
	c := vault.NewTestClient(map[string]string{"x": "1"})
	dests := []string{"secret/a", "secret/b", "secret/c"}
	results, err := Broadcast(c, "secret/src", dests, Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != len(dests) {
		t.Fatalf("expected %d results, got %d", len(dests), len(results))
	}
}

func TestDryRun_SourcePreservedInResult(t *testing.T) {
	c := vault.NewTestClient(map[string]string{"k": "v"})
	results, err := Broadcast(c, "secret/origin", []string{"secret/copy"}, Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Source != "secret/origin" {
		t.Errorf("expected source 'secret/origin', got %q", results[0].Source)
	}
}

func TestDryRun_BroadcastAtSet(t *testing.T) {
	c := vault.NewTestClient(map[string]string{"k": "v"})
	results, _ := Broadcast(c, "secret/src", []string{"secret/dst"}, Options{DryRun: true})
	if results[0].BroadcastAt.IsZero() {
		t.Error("expected BroadcastAt to be set")
	}
}
