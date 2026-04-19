package expire_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpatch/internal/expire"
)

func TestExpire_NilClientErrors(t *testing.T) {
	_, err := expire.Expire(nil, expire.Options{Paths: []string{"secret/a"}})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestExpire_EmptyPathsError(t *testing.T) {
	c := stubClient(map[string]map[string]interface{}{})
	_, err := expire.Expire(c, expire.Options{})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestExpire_DryRun_ResultFlagged(t *testing.T) {
	past := time.Now().Add(-24 * time.Hour).UTC().Format(time.RFC3339)
	c := stubClient(map[string]map[string]interface{}{
		"secret/old": {"_expire_at": past, "key": "value"},
	})
	res, err := expire.Expire(c, expire.Options{
		Paths:  []string{"secret/old"},
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if !res[0].DryRun {
		t.Error("expected DryRun to be true")
	}
	if res[0].Removed {
		t.Error("expected Removed to be false in dry-run")
	}
}

func TestExpire_NotYetExpired_Skipped(t *testing.T) {
	future := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	c := stubClient(map[string]map[string]interface{}{
		"secret/fresh": {"_expire_at": future},
	})
	res, err := expire.Expire(c, expire.Options{Paths: []string{"secret/fresh"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("expected no results, got %d", len(res))
	}
}

func TestExpire_NoExpiryKey_Skipped(t *testing.T) {
	c := stubClient(map[string]map[string]interface{}{
		"secret/plain": {"user": "admin"},
	})
	res, err := expire.Expire(c, expire.Options{Paths: []string{"secret/plain"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("expected no results, got %d", len(res))
	}
}
