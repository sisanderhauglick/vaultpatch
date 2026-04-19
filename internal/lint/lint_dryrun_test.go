package lint

import (
	"testing"
)

func TestLint_LintedAtSet(t *testing.T) {
	c := &stubClient{data: map[string]map[string]interface{}{
		"secret/x": {"k": "v"},
	}}
	results, err := Lint(c, Options{Paths: []string{"secret/x"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].LintedAt.IsZero() {
		t.Fatal("expected LintedAt to be set")
	}
}

func TestLint_PathPreservedInResult(t *testing.T) {
	c := &stubClient{data: map[string]map[string]interface{}{
		"secret/y": {"k": "v"},
	}}
	results, _ := Lint(c, Options{Paths: []string{"secret/y"}})
	if results[0].Path != "secret/y" {
		t.Fatalf("expected path %q, got %q", "secret/y", results[0].Path)
	}
}

func TestLint_MultiPath_AllReturned(t *testing.T) {
	c := &stubClient{data: map[string]map[string]interface{}{
		"secret/a": {"k": "v"},
		"secret/b": {"k": "v"},
	}}
	results, err := Lint(c, Options{Paths: []string{"secret/a", "secret/b"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestLint_DefaultRulesAppliedWhenNoneProvided(t *testing.T) {
	c := &stubClient{data: map[string]map[string]interface{}{
		"secret/z": {"BAD KEY": ""},
	}}
	results, _ := Lint(c, Options{Paths: []string{"secret/z"}})
	if len(results[0].Violations) < 2 {
		t.Fatalf("expected at least 2 violations from default rules, got %d", len(results[0].Violations))
	}
}
