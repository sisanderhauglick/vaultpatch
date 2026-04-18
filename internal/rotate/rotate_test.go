package rotate_test

import (
	"testing"

	"github.com/vaultpatch/vaultpatch/internal/rotate"
)

func TestRotate_NilClientErrors(t *testing.T) {
	_, err := rotate.Rotate(nil, rotate.Options{
		Path: "secret/data/app",
		Keys: []string{"password"},
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestRotate_EmptyPathErrors(t *testing.T) {
	_, err := rotate.Rotate(nil, rotate.Options{
		Path: "",
		Keys: []string{"password"},
	})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestRotate_EmptyKeysErrors(t *testing.T) {
	_, err := rotate.Rotate(nil, rotate.Options{
		Path: "secret/data/app",
		Keys: []string{},
	})
	if err == nil {
		t.Fatal("expected error for empty keys")
	}
}

func TestResult_DryRunFlagged(t *testing.T) {
	r := rotate.Result{
		Path:   "secret/data/app",
		Key:    "token",
		DryRun: true,
		Rotated: false,
	}
	if r.Rotated {
		t.Error("dry-run result should not be marked as rotated")
	}
	if !r.DryRun {
		t.Error("expected DryRun flag to be true")
	}
}

func TestResult_LiveRunFlagged(t *testing.T) {
	r := rotate.Result{
		Path:    "secret/data/app",
		Key:     "token",
		DryRun:  false,
		Rotated: true,
	}
	if !r.Rotated {
		t.Error("live result should be marked as rotated")
	}
}
