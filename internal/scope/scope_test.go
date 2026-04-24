package scope_test

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/scope"
)

// --- nil-client guard ---

func TestScope_NilClientErrors(t *testing.T) {
	_, err := scope.Scope(nil, []string{"secret/"}, scope.Options{})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestScope_EmptyRootsErrors(t *testing.T) {
	c := stubClient(map[string]map[string]string{})
	_, err := scope.Scope(c, []string{}, scope.Options{})
	if err == nil {
		t.Fatal("expected error for empty roots")
	}
}

// --- filterPaths (via Scope) ---

func TestScope_PrefixFiltersResults(t *testing.T) {
	data := map[string]map[string]string{
		"secret/app": {"key": "val"},
		"secret/db":  {"key": "val"},
	}
	c := stubClient(data)
	c.SetList("secret/", []string{"secret/app", "secret/db"})

	res, err := scope.Scope(c, []string{"secret/"}, scope.Options{Prefix: "secret/a"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if len(res[0].Paths) != 1 || res[0].Paths[0] != "secret/app" {
		t.Errorf("expected [secret/app], got %v", res[0].Paths)
	}
}

func TestScope_EmptyPrefix_ReturnsAll(t *testing.T) {
	c := stubClient(map[string]map[string]string{})
	c.SetList("secret/", []string{"secret/a", "secret/b", "secret/c"})

	res, err := scope.Scope(c, []string{"secret/"}, scope.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res[0].Paths) != 3 {
		t.Errorf("expected 3 paths, got %d", len(res[0].Paths))
	}
}

func TestScope_KeyFilter_ReturnsOnlyMatchingPaths(t *testing.T) {
	data := map[string]map[string]string{
		"secret/a": {"password": "s3cr3t"},
		"secret/b": {"token": "abc"},
	}
	c := stubClient(data)
	c.SetList("secret/", []string{"secret/a", "secret/b"})

	res, err := scope.Scope(c, []string{"secret/"}, scope.Options{Keys: []string{"password"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res[0].Paths) != 1 || res[0].Paths[0] != "secret/a" {
		t.Errorf("unexpected paths: %v", res[0].Paths)
	}
}

func TestScope_DryRun_FlaggedInResult(t *testing.T) {
	c := stubClient(map[string]map[string]string{})
	c.SetList("secret/", []string{"secret/x"})

	res, err := scope.Scope(c, []string{"secret/"}, scope.Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res[0].DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestScope_MatchedAt_Set(t *testing.T) {
	c := stubClient(map[string]map[string]string{})
	c.SetList("secret/", []string{})

	res, err := scope.Scope(c, []string{"secret/"}, scope.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res[0].MatchedAt.IsZero() {
		t.Error("expected MatchedAt to be set")
	}
}
