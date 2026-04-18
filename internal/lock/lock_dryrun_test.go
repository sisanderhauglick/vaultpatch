package lock_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpatch/internal/lock"
)

func TestDryRun_LockEntryHasExpiry(t *testing.T) {
	ttl := 15 * time.Minute
	res, err := lock.Lock(nil, "secret/app", "bot", ttl, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entry.ExpiresAt.IsZero() {
		t.Error("expected ExpiresAt to be set")
	}
	diff := res.Entry.ExpiresAt.Sub(res.Entry.LockedAt)
	if diff < ttl-time.Second || diff > ttl+time.Second {
		t.Errorf("expected expiry diff ~%s, got %s", ttl, diff)
	}
}

func TestDryRun_LockOwnerPreserved(t *testing.T) {
	res, err := lock.Lock(nil, "secret/app", "alice", time.Minute, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entry.Owner != "alice" {
		t.Errorf("expected owner %q, got %q", "alice", res.Entry.Owner)
	}
}

func TestDryRun_UnlockReleasedFlag(t *testing.T) {
	res, err := lock.Unlock(nil, "secret/app", "alice", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Released {
		t.Error("expected Released=true in dry-run unlock")
	}
	if !res.DryRun {
		t.Error("expected DryRun=true")
	}
}
