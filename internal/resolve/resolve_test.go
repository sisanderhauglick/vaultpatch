package resolve

import (
	"testing"
)

func TestResolve_NilClientErrors(t *testing.T) {
	_, err := Resolve(nil, Options{Paths: []string{"secret/app"}})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestResolve_EmptyPathsError(t *testing.T) {
	// We cannot construct a real vault.Client easily in unit tests,
	// so we test the guard before the client is used.
	_, err := Resolve(nil, Options{Paths: []string{}})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestParseRef_Valid(t *testing.T) {
	cases := []struct {
		input   string
		wantPath string
		wantKey  string
	}{
		{"ref:secret/db#password", "secret/db", "password"},
		{"ref:kv/prod/api#token", "kv/prod/api", "token"},
	}
	for _, tc := range cases {
		path, key, err := parseRef(tc.input)
		if err != nil {
			t.Fatalf("parseRef(%q) unexpected error: %v", tc.input, err)
		}
		if path != tc.wantPath {
			t.Errorf("path: got %q, want %q", path, tc.wantPath)
		}
		if key != tc.wantKey {
			t.Errorf("key: got %q, want %q", key, tc.wantKey)
		}
	}
}

func TestParseRef_Invalid(t *testing.T) {
	cases := []string{
		"ref:",
		"ref:nohash",
		"ref:#onlykey",
		"ref:path#",
		"notaref",
	}
	for _, tc := range cases {
		_, _, err := parseRef(tc)
		if err == nil {
			t.Errorf("parseRef(%q) expected error, got nil", tc)
		}
	}
}

func TestResult_DryRunFlagged(t *testing.T) {
	r := Result{
		Path:   "secret/app",
		DryRun: true,
		Resolved: map[string]string{
			"DB_PASS": "s3cr3t",
		},
	}
	if !r.DryRun {
		t.Error("expected DryRun to be true")
	}
	if r.Resolved["DB_PASS"] != "s3cr3t" {
		t.Errorf("unexpected resolved value: %v", r.Resolved["DB_PASS"])
	}
}

func TestResult_UnresolvedTracked(t *testing.T) {
	r := Result{
		Path:       "secret/app",
		Unresolved: []string{"MISSING_KEY"},
	}
	if len(r.Unresolved) != 1 || r.Unresolved[0] != "MISSING_KEY" {
		t.Errorf("unexpected unresolved list: %v", r.Unresolved)
	}
}

func TestResult_ResolvedAtSet(t *testing.T) {
	r := Result{}
	if !r.ResolvedAt.IsZero() {
		t.Error("expected zero ResolvedAt for empty result")
	}
}
