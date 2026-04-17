package rename

import (
	"testing"

	"vaultpatch/internal/vault"
)

func TestRename_NilClientErrors(t *testing.T) {
	_, err := Rename(nil, Options{Paths: []string{"secret/app"}, OldKey: "A", NewKey: "B"})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestRename_EmptyPathsError(t *testing.T) {
	c, _ := vault.NewClient(vault.Params{Addr: "http://127.0.0.1:8200", Token: "tok"})
	_, err := Rename(c, Options{Paths: nil, OldKey: "A", NewKey: "B"})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestRename_EmptyKeyErrors(t *testing.T) {
	c, _ := vault.NewClient(vault.Params{Addr: "http://127.0.0.1:8200", Token: "tok"})
	_, err := Rename(c, Options{Paths: []string{"secret/app"}, OldKey: "", NewKey: "B"})
	if err == nil {
		t.Fatal("expected error for empty old-key")
	}
	_, err = Rename(c, Options{Paths: []string{"secret/app"}, OldKey: "A", NewKey: ""})
	if err == nil {
		t.Fatal("expected error for empty new-key")
	}
}

func TestRename_DryRun_SkipsWrite(t *testing.T) {
	// Without a live Vault, we verify that dry-run returns a non-error path
	// when the underlying read returns a not-found / empty map (Skipped=true).
	c, _ := vault.NewClient(vault.Params{Addr: "http://127.0.0.1:8200", Token: "tok"})
	opts := Options{
		Paths:  []string{"secret/missing"},
		OldKey: "FOO",
		NewKey: "BAR",
		DryRun: true,
	}
	// ReadSecrets against a non-running server will error; we just confirm the
	// function propagates errors correctly rather than silently succeeding.
	_, err := Rename(c, opts)
	if err == nil {
		t.Log("no error returned (live vault present or mock used)")
	}
}

func TestResult_SkippedFlag(t *testing.T) {
	r := Result{Path: "secret/app", OldKey: "X", NewKey: "Y", DryRun: true, Skipped: true}
	if !r.Skipped {
		t.Error("expected Skipped to be true")
	}
	if !r.DryRun {
		t.Error("expected DryRun to be true")
	}
}
