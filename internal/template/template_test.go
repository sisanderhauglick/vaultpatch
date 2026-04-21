package template

import (
	"testing"
)

func TestRender_StaticPath_NoChange(t *testing.T) {
	res, err := Render("secret/app/config", map[string]string{"key": "value"}, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Path != "secret/app/config" {
		t.Errorf("expected unchanged path, got %q", res.Path)
	}
	if res.Rendered["key"] != "value" {
		t.Errorf("expected unchanged value, got %q", res.Rendered["key"])
	}
}

func TestRender_EmptyPathErrors(t *testing.T) {
	_, err := Render("", map[string]string{}, Options{})
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestRender_DynamicPath(t *testing.T) {
	vars := map[string]string{"env": "staging"}
	res, err := Render("secret/{{.env}}/db", map[string]string{}, Options{Vars: vars})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Path != "secret/staging/db" {
		t.Errorf("expected rendered path, got %q", res.Path)
	}
}

func TestRender_DynamicValue(t *testing.T) {
	vars := map[string]string{"region": "us-east-1"}
	secrets := map[string]string{"endpoint": "https://{{.region}}.example.com"}
	res, err := Render("secret/app", secrets, Options{Vars: vars})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Rendered["endpoint"] != "https://us-east-1.example.com" {
		t.Errorf("unexpected rendered value: %q", res.Rendered["endpoint"])
	}
}

func TestRender_MissingVar_NonStrict_EmptyString(t *testing.T) {
	secrets := map[string]string{"val": "prefix-{{.missing}}-suffix"}
	res, err := Render("secret/app", secrets, Options{})
	if err != nil {
		t.Fatalf("unexpected error in non-strict mode: %v", err)
	}
	if res.Rendered["val"] != "prefix--suffix" {
		t.Errorf("expected empty substitution, got %q", res.Rendered["val"])
	}
}

func TestRender_MissingVar_Strict_Errors(t *testing.T) {
	secrets := map[string]string{"val": "{{.missing}}"}
	_, err := Render("secret/app", secrets, Options{Strict: true})
	if err == nil {
		t.Fatal("expected error in strict mode for missing variable, got nil")
	}
}

func TestRender_NilVars_DefaultsToEmpty(t *testing.T) {
	res, err := Render("secret/app", map[string]string{"k": "v"}, Options{Vars: nil})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Rendered["k"] != "v" {
		t.Errorf("expected value preserved, got %q", res.Rendered["k"])
	}
}
