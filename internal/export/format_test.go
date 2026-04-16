package export

import (
	"testing"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []string{"json", "yaml", "env"}
	for _, c := range cases {
		f, err := ParseFormat(c)
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", c, err)
		}
		if string(f) != c {
			t.Errorf("ParseFormat(%q) = %q, want %q", c, f, c)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := ParseFormat("xml")
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestSupportedFormats(t *testing.T) {
	formats := SupportedFormats()
	if len(formats) != 3 {
		t.Errorf("expected 3 supported formats, got %d", len(formats))
	}
}
