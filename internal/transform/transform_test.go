package transform_test

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/transform"
)

func TestTransform_NilClientErrors(t *testing.T) {
	_, err := transform.Transform(nil, transform.Options{
		Paths: []string{"secret/data/app"},
		Rules: []transform.Rule{{Op: transform.OpUppercase}},
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestTransform_EmptyPathsError(t *testing.T) {
	c := stubClient()
	_, err := transform.Transform(c, transform.Options{
		Paths: nil,
		Rules: []transform.Rule{{Op: transform.OpUppercase}},
	})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestTransform_EmptyRulesError(t *testing.T) {
	c := stubClient()
	_, err := transform.Transform(c, transform.Options{
		Paths: []string{"secret/data/app"},
		Rules: nil,
	})
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestApplyRules_Uppercase(t *testing.T) {
	data := map[string]string{"key": "hello"}
	applyRulesExported(data, []transform.Rule{{Op: transform.OpUppercase}})
	if data["key"] != "HELLO" {
		t.Fatalf("expected HELLO, got %s", data["key"])
	}
}

func TestApplyRules_Lowercase(t *testing.T) {
	data := map[string]string{"key": "WORLD"}
	applyRulesExported(data, []transform.Rule{{Op: transform.OpLowercase}})
	if data["key"] != "world" {
		t.Fatalf("expected world, got %s", data["key"])
	}
}

func TestApplyRules_TrimSpace(t *testing.T) {
	data := map[string]string{"key": "  spaced  "}
	applyRulesExported(data, []transform.Rule{{Op: transform.OpTrimSpace}})
	if data["key"] != "spaced" {
		t.Fatalf("expected 'spaced', got %q", data["key"])
	}
}

func TestApplyRules_Prefix(t *testing.T) {
	data := map[string]string{"key": "value"}
	applyRulesExported(data, []transform.Rule{{Op: transform.OpPrefix, Value: "pre_"}})
	if data["key"] != "pre_value" {
		t.Fatalf("expected pre_value, got %s", data["key"])
	}
}

func TestApplyRules_Suffix(t *testing.T) {
	data := map[string]string{"key": "value"}
	applyRulesExported(data, []transform.Rule{{Op: transform.OpSuffix, Value: "_suf"}})
	if data["key"] != "value_suf" {
		t.Fatalf("expected value_suf, got %s", data["key"])
	}
}

func TestApplyRules_SubsetKeys_OnlyTargeted(t *testing.T) {
	data := map[string]string{"a": "hello", "b": "world"}
	applyRulesExported(data, []transform.Rule{
		{Op: transform.OpUppercase, Keys: []string{"a"}},
	})
	if data["a"] != "HELLO" {
		t.Fatalf("expected HELLO for key a, got %s", data["a"])
	}
	if data["b"] != "world" {
		t.Fatalf("expected world unchanged for key b, got %s", data["b"])
	}
}

func TestResult_DryRunFlagged(t *testing.T) {
	c := stubClient()
	res, err := transform.Transform(c, transform.Options{
		Paths:  []string{"secret/data/app"},
		Rules:  []transform.Rule{{Op: transform.OpUppercase}},
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) == 0 {
		t.Fatal("expected at least one result")
	}
	if !res[0].DryRun {
		t.Fatal("expected DryRun to be true")
	}
}
