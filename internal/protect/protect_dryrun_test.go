package protect_test

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/protect"
	"github.com/your-org/vaultpatch/internal/vault"
)

// stubClient returns a minimal non-nil *vault.Client for unit tests that
// exercise dry-run paths (no real Vault calls are made).
func stubClient() *vault.Client {
	c, _ := vault.NewClient(vault.Config{
		Address: "http://127.0.0.1:8200",
		Token:   "root",
	})
	return c
}

func TestDryRun_ProtectResultFlagged(t *testing.T) {
	results, err := protect.Protect(stubClient(), protect.Options{
		Paths:  []string{"secret/app"},
		Owner:  "alice",
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
	if !results[0].Protected {
		t.Error("expected Protected to be true in dry-run")
	}
}

func TestDryRun_UnprotectResultFlagged(t *testing.T) {
	results, err := protect.Unprotect(stubClient(), protect.Options{
		Paths:  []string{"secret/app"},
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
	if results[0].Protected {
		t.Error("expected Protected to be false after unprotect")
	}
}

func TestDryRun_OwnerPreservedInResult(t *testing.T) {
	results, _ := protect.Protect(stubClient(), protect.Options{
		Paths:  []string{"secret/app"},
		Owner:  "bob",
		DryRun: true,
	})
	if results[0].Owner != "bob" {
		t.Errorf("expected owner %q, got %q", "bob", results[0].Owner)
	}
}

func TestDryRun_ProtectedAtSet(t *testing.T) {
	results, _ := protect.Protect(stubClient(), protect.Options{
		Paths:  []string{"secret/app"},
		Owner:  "alice",
		DryRun: true,
	})
	if results[0].ProtectedAt.IsZero() {
		t.Error("expected ProtectedAt to be set")
	}
}
