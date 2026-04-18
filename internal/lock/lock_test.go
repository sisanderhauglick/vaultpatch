package lock_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpatch/internal/lock"
)

func TestLock_NilClientErrors(t *testing.T) {
	_, err := lock.Lock(nil, "secret/foo", "ci", time.Minute, false)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestLock_EmptyPathErrors(t *testing.T) {
	_, err := lock.Lock(nil, "", "ci", time.Minute, false)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestLock_EmptyOwnerErrors(t *testing.T) {
	_, err := lock.Lock(nil, "secret/foo", "", time.Minute, false)
	if err == nil {
		t.Fatal("expected error for empty owner")
	}
}

func TestUnlock_NilClientErrors(t *testing.T) {
	_, err := lock.Unlock(nil, "secret/foo", "ci", false)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestUnlock_EmptyPathErrors(t *testing.T) {
	_, err := lock.Unlock(nil, "", "ci", false)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestLock_DryRun_DoesNotWrite(t *testing.T) {
	res, err := lock.Lock(nil, "secret/foo", "ci", time.Minute, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.DryRun {
		t.Error("expected DryRun to be true")
	}
	if !res.Acquired {
		t.Error("expected Acquired to be true in dry-run")
	}
	if res.Entry == nil {
		t.Error("expected Entry to be populated in dry-run")
	}
}

func TestUnlock_DryRun_DoesNotDelete(t *testing.T) {
	res, err := lock.Unlock(nil, "secret/foo", "ci", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.DryRun {
		t.Error("expected DryRun to be true")
	}
	if !res.Released {
		t.Error("expected Released to be true in dry-run")
	}
}

func TestResult_PathPreserved(t *testing.T) {
	res, _ := lock.Lock(nil, "secret/myapp", "dev", 5*time.Minute, true)
	if res.Path != "secret/myapp" {
		t.Errorf("expected path %q, got %q", "secret/myapp", res.Path)
	}
}
