package format

import (
	"strings"
	"testing"
)

func TestParseStyle_Valid(t *testing.T) {
	for _, s := range SupportedStyles() {
		_, err := ParseStyle(s)
		if err != nil {
			t.Errorf("expected %q to be valid, got error: %v", s, err)
		}
	}
}

func TestParseStyle_Invalid(t *testing.T) {
	_, err := ParseStyle("xml")
	if err == nil {
		t.Fatal("expected error for unsupported style")
	}
	if !strings.Contains(err.Error(), "xml") {
		t.Errorf("error should mention invalid style, got: %v", err)
	}
}

func TestSupportedStyles(t *testing.T) {
	styles := SupportedStyles()
	if len(styles) != 3 {
		t.Errorf("expected 3 styles, got %d", len(styles))
	}
}

func TestRender_Table(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out := Render(secrets, StyleTable)
	if !strings.Contains(out, "KEY") || !strings.Contains(out, "VALUE") {
		t.Error("table header missing")
	}
	if !strings.Contains(out, "FOO") || !strings.Contains(out, "bar") {
		t.Error("table row missing expected content")
	}
}

func TestRender_List(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	out := Render(secrets, StyleList)
	if !strings.Contains(out, "A=1") {
		t.Errorf("expected A=1 in list output, got: %s", out)
	}
}

func TestRender_CSV(t *testing.T) {
	secrets := map[string]string{"X": "y"}
	out := Render(secrets, StyleCSV)
	if !strings.Contains(out, "key,value") {
		t.Error("csv header missing")
	}
	if !strings.Contains(out, "X,y") {
		t.Errorf("expected X,y in csv output, got: %s", out)
	}
}

func TestRender_Empty(t *testing.T) {
	out := Render(map[string]string{}, StyleList)
	if out != "" {
		t.Errorf("expected empty output for empty secrets, got: %q", out)
	}
}
