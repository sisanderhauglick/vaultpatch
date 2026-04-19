package inspect_test

import (
	"testing"

	"github.com/youorg/vaultpatch/internal/inspect"
)

func TestInspect_NilClientErrors(t *testing.T) {
	_, err := inspect.Inspect(nil, inspect.Options{Paths: []string{"secret/app"}})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestInspect_EmptyPathsError(t *testing.T) {
	_, err := inspect.Inspect(nil, inspect.Options{})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestResult_KeyCountMatchesKeys(t *testing.T) {
	r := inspect.Result{
		Path:     "secret/app",
		Keys:     []string{"DB_HOST", "DB_PASS", "API_KEY"},
		KeyCount: 3,
	}
	if r.KeyCount != len(r.Keys) {
		t.Errorf("KeyCount %d does not match len(Keys) %d", r.KeyCount, len(r.Keys))
	}
}

func TestResult_FetchedAtSet(t *testing.T) {
	r := inspect.Result{}
	if !r.FetchedAt.IsZero() {
		t.Error("expected zero FetchedAt on empty result")
	}
}

func TestResult_PathPreserved(t *testing.T) {
	path := "secret/service/config"
	r := inspect.Result{Path: path}
	if r.Path != path {
		t.Errorf("expected path %q, got %q", path, r.Path)
	}
}
