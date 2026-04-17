package validate

import (
	"testing"
)

func TestValidate_RequiredKeyMissing(t *testing.T) {
	rules := []Rule{{Key: "API_KEY", Required: true}}
	v := Validate("secret/app", map[string]string{}, rules)
	if len(v) != 1 || v[0].Message != "required key is missing" {
		t.Fatalf("expected required-key violation, got %+v", v)
	}
}

func TestValidate_OptionalKeyAbsent_NoViolation(t *testing.T) {
	rules := []Rule{{Key: "OPTIONAL", Required: false}}
	v := Validate("secret/app", map[string]string{}, rules)
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %+v", v)
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	rules := []Rule{{Key: "DB_PASS", Required: true}}
	v := Validate("secret/db", map[string]string{"DB_PASS": "   "}, rules)
	if len(v) != 1 {
		t.Fatalf("expected empty-value violation, got %+v", v)
	}
}

func TestValidate_PatternMatch_Pass(t *testing.T) {
	rules := []Rule{{Key: "PORT", Pattern: `^\d+$`}}
	v := Validate("secret/svc", map[string]string{"PORT": "8080"}, rules)
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %+v", v)
	}
}

func TestValidate_PatternMatch_Fail(t *testing.T) {
	rules := []Rule{{Key: "PORT", Pattern: `^\d+$`}}
	v := Validate("secret/svc", map[string]string{"PORT": "not-a-port"}, rules)
	if len(v) != 1 {
		t.Fatalf("expected pattern violation, got %+v", v)
	}
}

func TestValidate_InvalidPattern(t *testing.T) {
	rules := []Rule{{Key: "X", Pattern: `[invalid`}}
	v := Validate("secret/x", map[string]string{"X": "value"}, rules)
	if len(v) != 1 {
		t.Fatalf("expected invalid-pattern violation, got %+v", v)
	}
}

func TestValidate_NoRules_NoViolations(t *testing.T) {
	v := Validate("secret/any", map[string]string{"FOO": "bar"}, nil)
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %+v", v)
	}
}

func TestViolation_Error(t *testing.T) {
	v := Violation{Path: "secret/app", Key: "KEY", Message: "required key is missing"}
	want := "secret/app/KEY: required key is missing"
	if v.Error() != want {
		t.Fatalf("expected %q got %q", want, v.Error())
	}
}
