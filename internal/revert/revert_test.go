package revert_test

import (
	"testing"

	"github.com/user/vaultpatch/internal/revert"
	"github.com/user/vaultpatch/internal/vault"
)

func nilClient() *vault.Client { return nil }

func TestRevert_NilClientErrors(t *testing.T) {
	_, err := revert.Revert(nil, revert.Options{
		Paths:  []string{"secret/app"},
		Before: map[string]string{"KEY": "old"},
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestRevert_EmptyPathsError(t *testing.T) {
	c, _ := vault.NewClient(vault.Options{Addr: "http://127.0.0.1:8200", Token: "t"})
	_, err := revert.Revert(c, revert.Options{
		Before: map[string]string{"KEY": "old"},
	})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestRevert_EmptyBeforeError(t *testing.T) {
	c, _ := vault.NewClient(vault.Options{Addr: "http://127.0.0.1:8200", Token: "t"})
	_, err := revert.Revert(c, revert.Options{
		Paths: []string{"secret/app"},
	})
	if err == nil {
		t.Fatal("expected error for empty before map")
	}
}

func TestResult_DryRunFlagged(t *testing.T) {
	res := revert.Result{DryRun: true, Reverted: []string{"KEY"}}
	if !res.DryRun {
		t.Error("expected DryRun to be true")
	}
	if len(res.Reverted) != 1 || res.Reverted[0] != "KEY" {
		t.Errorf("unexpected reverted keys: %v", res.Reverted)
	}
}

func TestResult_SkippedWhenValueUnchanged(t *testing.T) {
	res := revert.Result{
		Path:    "secret/app",
		Skipped: []string{"UNCHANGED"},
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped key, got %d", len(res.Skipped))
	}
}

func TestResult_PathPreserved(t *testing.T) {
	res := revert.Result{Path: "secret/staging"}
	if res.Path != "secret/staging" {
		t.Errorf("unexpected path: %s", res.Path)
	}
}
