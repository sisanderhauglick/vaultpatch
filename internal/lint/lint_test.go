package lint

import (
	"errors"
	"testing"
)

type stubClient struct {
	data map[string]map[string]interface{}
	err  error
}

func (s *stubClient) Read(path string) (map[string]interface{}, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.data[path], nil
}

func TestLint_NilClientErrors(t *testing.T) {
	_, err := Lint(nil, Options{Paths: []string{"secret/a"}})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestLint_EmptyPathsError(t *testing.T) {
	_, err := Lint(&stubClient{}, Options{})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestLint_ReadError(t *testing.T) {
	c := &stubClient{err: errors.New("vault down")}
	_, err := Lint(c, Options{Paths: []string{"secret/a"}})
	if err == nil {
		t.Fatal("expected error on read failure")
	}
}

func TestLint_NoViolations(t *testing.T) {
	c := &stubClient{data: map[string]map[string]interface{}{
		"secret/a": {"api_key": "abc123"},
	}}
	results, err := Lint(c, Options{Paths: []string{"secret/a"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results[0].Violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(results[0].Violations))
	}
}

func TestLint_EmptyValueViolation(t *testing.T) {
	c := &stubClient{data: map[string]map[string]interface{}{
		"secret/a": {"token": ""},
	}}
	results, err := Lint(c, Options{Paths: []string{"secret/a"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results[0].Violations) == 0 {
		t.Fatal("expected at least one violation for empty value")
	}
}

func TestLint_UppercaseKeyViolation(t *testing.T) {
	c := &stubClient{data: map[string]map[string]interface{}{
		"secret/b": {"API_KEY": "val"},
	}}
	results, _ := Lint(c, Options{Paths: []string{"secret/b"}})
	found := false
	for _, v := range results[0].Violations {
		if v.Rule == "no-uppercase-key" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected no-uppercase-key violation")
	}
}

func TestLint_CustomRule(t *testing.T) {
	c := &stubClient{data: map[string]map[string]interface{}{
		"secret/c": {"key": "forbidden"},
	}}
	custom := Rule{
		Name: "no-forbidden", Message: "value is forbidden",
		Check: func(_, v string) bool { return v == "forbidden" },
	}
	results, _ := Lint(c, Options{Paths: []string{"secret/c"}, Rules: []Rule{custom}})
	if len(results[0].Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(results[0].Violations))
	}
}
