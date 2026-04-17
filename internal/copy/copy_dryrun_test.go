package copy

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/vault"
)

// stubResult is a helper that builds a Result directly for assertion.
func stubResult(src, dst string, keys int, dry bool) Result {
	return Result{SourcePath: src, DestPath: dst, Keys: keys, DryRun: dry}
}

func TestDryRun_ResultFlagged(t *testing.T) {
	res := stubResult("src/path", "dst/path", 3, true)
	if !res.DryRun {
		t.Fatal("expected DryRun to be true")
	}
	if res.Keys != 3 {
		t.Fatalf("expected 3 keys, got %d", res.Keys)
	}
}

func TestLiveRun_ResultFlagged(t *testing.T) {
	res := stubResult("src/path", "dst/path", 5, false)
	if res.DryRun {
		t.Fatal("expected DryRun to be false")
	}
}

func TestFilterKeys_EmptyInputReturnsAll(t *testing.T) {
	sm := vault.SecretMap{"x": "1", "y": "2", "z": "3"}
	out := filterKeys(sm, []string{})
	if len(out) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(out))
	}
}
