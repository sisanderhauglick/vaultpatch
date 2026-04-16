package export

import (
	"bytes"
	"strings"
	"testing"
)

var sampleSecrets = map[string]string{
	"db_password": "s3cr3t",
	"api_key":     "abc123",
}

func TestExport_JSON(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(sampleSecrets, FormatJSON, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "db_password") {
		t.Errorf("expected db_password in JSON output, got: %s", out)
	}
	if !strings.Contains(out, "s3cr3t") {
		t.Errorf("expected s3cr3t in JSON output, got: %s", out)
	}
}

func TestExport_YAML(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(sampleSecrets, FormatYAML, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "api_key") {
		t.Errorf("expected api_key in YAML output, got: %s", out)
	}
}

func TestExport_Env(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(sampleSecrets, FormatEnv, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Errorf("expected DB_PASSWORD in env output, got: %s", out)
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Export(sampleSecrets, Format("toml"), &buf)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestExport_EmptySecrets(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(map[string]string{}, FormatJSON, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output even for empty secrets")
	}
}
