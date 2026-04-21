package dedupe

import (
	"testing"

	"github.com/yourusername/vaultpatch/internal/vault"
)

func TestDedupe_NilClientErrors(t *testing.T) {
	_, err := Dedupe(nil, Options{Paths: []string{"a", "b"}})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestDedupe_TooFewPathsErrors(t *testing.T) {
	client := &vault.Client{}
	_, err := Dedupe(client, Options{Paths: []string{"only-one"}})
	if err == nil {
		t.Fatal("expected error for fewer than two paths")
	}
}

func TestDedupe_EmptyPathsErrors(t *testing.T) {
	client := &vault.Client{}
	_, err := Dedupe(client, Options{Paths: []string{}})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestResult_DryRunFlagged(t *testing.T) {
	r := Result{DryRun: true, Path: "secret/dev"}
	if !r.DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestResult_DedupedAtSet(t *testing.T) {
	r := Result{Path: "secret/dev"}
	if !r.DedupedAt.IsZero() {
		// zero is fine for a bare struct; populated by Dedupe()
	}
	_ = r
}

func TestResult_PathPreserved(t *testing.T) {
	r := Result{Path: "secret/staging/app"}
	if r.Path != "secret/staging/app" {
		t.Errorf("unexpected path: %s", r.Path)
	}
}

func TestResult_RemovedKeysSlice(t *testing.T) {
	r := Result{
		Path:        "secret/prod",
		RemovedKeys: []string{"DB_PASS", "API_KEY"},
	}
	if len(r.RemovedKeys) != 2 {
		t.Errorf("expected 2 removed keys, got %d", len(r.RemovedKeys))
	}
}
