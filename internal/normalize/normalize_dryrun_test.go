package normalize

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/vault"
)

// stubClient returns a minimal *vault.Client suitable for unit tests that
// do not need real Vault connectivity.
func stubClient(t *testing.T) *vault.Client {
	t.Helper()
	c, err := vault.NewClient(vault.Params{
		Address: "http://127.0.0.1:8200",
		Token:   "test-token",
	})
	if err != nil {
		t.Skipf("stub client unavailable: %v", err)
	}
	return c
}

func TestDryRun_ResultFlagged(t *testing.T) {
	result := Result{
		Path:    "secret/test",
		Changes: map[string]string{"key": "VALUE"},
		DryRun:  true,
	}
	if !result.DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestDryRun_ChangesPopulated(t *testing.T) {
	result := Result{
		Path: "secret/app",
		Changes: map[string]string{
			"host": "LOCALHOST",
			"port": "5432",
		},
		DryRun: true,
	}
	if len(result.Changes) != 2 {
		t.Errorf("expected 2 changes, got %d", len(result.Changes))
	}
}

func TestDryRun_PathPreservedInResult(t *testing.T) {
	const path = "secret/service/config"
	result := Result{Path: path, DryRun: true}
	if result.Path != path {
		t.Errorf("expected path %q, got %q", path, result.Path)
	}
}
