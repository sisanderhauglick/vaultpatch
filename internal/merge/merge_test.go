package merge

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/vault"
)

func TestMerge_NilClientErrors(t *testing.T) {
	_, err := Merge(nil, Options{Sources: []string{"a"}, Destination: "dst"})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestMerge_EmptySourcesErrors(t *testing.T) {
	client := &vault.Client{}
	_, err := Merge(client, Options{Sources: []string{}, Destination: "dst"})
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestMerge_EmptyDestinationErrors(t *testing.T) {
	client := &vault.Client{}
	_, err := Merge(client, Options{Sources: []string{"a"}, Destination: ""})
	if err == nil {
		t.Fatal("expected error for empty destination")
	}
}

func TestContainsKey_Found(t *testing.T) {
	if !containsKey([]string{"foo", "bar"}, "bar") {
		t.Fatal("expected bar to be found")
	}
}

func TestContainsKey_NotFound(t *testing.T) {
	if containsKey([]string{"foo"}, "baz") {
		t.Fatal("expected baz not to be found")
	}
}

func TestContainsKey_EmptyList(t *testing.T) {
	if containsKey([]string{}, "any") {
		t.Fatal("expected no match on empty list")
	}
}

func TestResult_DryRunFlagged(t *testing.T) {
	r := Result{Destination: "dst", Merged: 3, DryRun: true}
	if !r.DryRun {
		t.Fatal("expected DryRun to be true")
	}
	if r.Merged != 3 {
		t.Fatalf("expected Merged=3, got %d", r.Merged)
	}
}
