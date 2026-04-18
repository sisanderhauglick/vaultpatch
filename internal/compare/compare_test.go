package compare

import (
	"testing"

	"github.com/vaultpatch/internal/vault"
)

func TestCompare_NilClientErrors(t *testing.T) {
	_, err := Compare(nil, "secret/a", "secret/b")
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestCompare_EmptySourceErrors(t *testing.T) {
	c := &vault.Client{}
	_, err := Compare(c, "", "secret/b")
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestCompare_EmptyDestErrors(t *testing.T) {
	c := &vault.Client{}
	_, err := Compare(c, "secret/a", "")
	if err == nil {
		t.Fatal("expected error for empty destination")
	}
}

func TestResult_OnlyInSrc(t *testing.T) {
	res := &Result{
		SourcePath: "secret/a",
		DestPath:   "secret/b",
		OnlyInSrc:  map[string]string{"key1": "val1"},
		OnlyInDst:  map[string]string{},
		Differ:     map[string][2]string{},
		Match:      map[string]string{},
	}
	if _, ok := res.OnlyInSrc["key1"]; !ok {
		t.Fatal("expected key1 in OnlyInSrc")
	}
}

func TestResult_Differ(t *testing.T) {
	res := &Result{
		Differ: map[string][2]string{
			"token": {"abc", "xyz"},
		},
		OnlyInSrc: map[string]string{},
		OnlyInDst: map[string]string{},
		Match:     map[string]string{},
	}
	pair, ok := res.Differ["token"]
	if !ok {
		t.Fatal("expected token in Differ")
	}
	if pair[0] != "abc" || pair[1] != "xyz" {
		t.Fatalf("unexpected values: %v", pair)
	}
}

func TestResult_Match(t *testing.T) {
	res := &Result{
		Match:     map[string]string{"shared": "same"},
		OnlyInSrc: map[string]string{},
		OnlyInDst: map[string]string{},
		Differ:    map[string][2]string{},
	}
	if v := res.Match["shared"]; v != "same" {
		t.Fatalf("expected 'same', got %q", v)
	}
}
