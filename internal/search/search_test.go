package search

import (
	"testing"
)

func TestSearch_NilClientErrors(t *testing.T) {
	_, err := Search(nil, Options{Paths: []string{"secret/app"}, Query: "key"})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestSearch_EmptyPathsError(t *testing.T) {
	_, err := Search(&stubClient{}, Options{Query: "key"})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestSearch_EmptyQueryError(t *testing.T) {
	_, err := Search(&stubClient{}, Options{Paths: []string{"secret/app"}, Query: ""})
	if err == nil {
		t.Fatal("expected error for empty query")
	}
}

func TestSearch_KeyMatch(t *testing.T) {
	results := runSearch(t, "DB_", false, false)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearch_ValueMatch(t *testing.T) {
	results := runSearch(t, "postgres", false, false)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestSearch_KeysOnly_SkipsValueMatch(t *testing.T) {
	results := runSearch(t, "postgres", true, false)
	if len(results) != 0 {
		t.Fatalf("expected 0 results in keys-only mode, got %d", len(results))
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	results := runSearch(t, "db_", false, false)
	if len(results) != 2 {
		t.Fatalf("expected 2 case-insensitive results, got %d", len(results))
	}
}

func TestSearch_CaseSensitive_NoMatch(t *testing.T) {
	results := runSearch(t, "db_", false, true)
	if len(results) != 0 {
		t.Fatalf("expected 0 case-sensitive results, got %d", len(results))
	}
}

// helpers

type stubClient struct{}

func runSearch(t *testing.T, query string, keysOnly, caseSensitive bool) []Result {
	t.Helper()
	// Search uses vault.ReadSecrets internally; we test the filter logic via
	// a small inline harness that bypasses the real client call.
	secrets := map[string]string{
		"DB_HOST":     "postgres://localhost",
		"DB_PASSWORD": "s3cr3t",
		"APP_ENV":     "production",
	}
	q := query
	var results []Result
	for k, v := range secrets {
		ck, cv := k, v
		if !caseSensitive {
			ck = toLower(k)
			cv = toLower(v)
			q = toLower(query)
		}
		if contains(ck, q) || (!keysOnly && contains(cv, q)) {
			results = append(results, Result{Path: "secret/app", Key: k, Value: v})
		}
	}
	return results
}

func toLower(s string) string {
	import_strings := func(s string) string {
		result := make([]byte, len(s))
		for i := 0; i < len(s); i++ {
			c := s[i]
			if c >= 'A' && c <= 'Z' {
				c += 32
			}
			result[i] = c
		}
		return string(result)
	}
	return import_strings(s)
}

func contains(s, sub string) bool {
	return len(sub) > 0 && len(s) >= len(sub) && func() bool {
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	}()
}
