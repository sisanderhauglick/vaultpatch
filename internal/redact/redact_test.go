package redact

import (
	"testing"
)

func secrets() map[string]map[string]string {
	return map[string]map[string]string{
		"secret/app": {
			"password": "s3cr3t",
			"host":     "localhost",
		},
	}
}

func TestRedact_AllKeys_WhenNoneSpecified(t *testing.T) {
	results := Redact(secrets(), Options{})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	for _, v := range results[0].Redacted {
		if v != "[REDACTED]" {
			t.Errorf("expected [REDACTED], got %q", v)
		}
	}
}

func TestRedact_SubsetKeys(t *testing.T) {
	results := Redact(secrets(), Options{Keys: []string{"password"}})
	r := results[0].Redacted
	if r["password"] != "[REDACTED]" {
		t.Errorf("password should be redacted")
	}
	if r["host"] != "localhost" {
		t.Errorf("host should be preserved, got %q", r["host"])
	}
}

func TestRedact_CustomReplacement(t *testing.T) {
	results := Redact(secrets(), Options{Replacement: "***"})
	for _, v := range results[0].Redacted {
		if v != "***" {
			t.Errorf("expected ***, got %q", v)
		}
	}
}

func TestRedact_CaseInsensitiveKey(t *testing.T) {
	results := Redact(secrets(), Options{Keys: []string{"PASSWORD"}})
	if results[0].Redacted["password"] != "[REDACTED]" {
		t.Errorf("case-insensitive match failed")
	}
}

func TestRedact_OriginalPreserved(t *testing.T) {
	results := Redact(secrets(), Options{})
	if results[0].Original["password"] != "s3cr3t" {
		t.Errorf("original should be preserved")
	}
}

func TestRedact_EmptySecrets(t *testing.T) {
	results := Redact(map[string]map[string]string{}, Options{})
	if len(results) != 0 {
		t.Errorf("expected empty results")
	}
}
