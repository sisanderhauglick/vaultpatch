package tag

import (
	"testing"
)

func TestParseTags_Empty(t *testing.T) {
	got := parseTags("")
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestParseTags_Single(t *testing.T) {
	got := parseTags("prod")
	if len(got) != 1 || got[0] != "prod" {
		t.Fatalf("unexpected tags: %v", got)
	}
}

func TestParseTags_Multiple(t *testing.T) {
	got := parseTags("prod,staging,dev")
	if len(got) != 3 {
		t.Fatalf("expected 3 tags, got %v", got)
	}
}

func TestApplyChanges_Add(t *testing.T) {
	out := applyChanges([]string{"prod"}, []string{"staging"}, nil)
	if !containsTag(out, "prod") || !containsTag(out, "staging") {
		t.Fatalf("expected prod and staging, got %v", out)
	}
}

func TestApplyChanges_Remove(t *testing.T) {
	out := applyChanges([]string{"prod", "staging"}, nil, []string{"staging"})
	if containsTag(out, "staging") {
		t.Fatalf("staging should have been removed, got %v", out)
	}
	if !containsTag(out, "prod") {
		t.Fatalf("prod should remain, got %v", out)
	}
}

func TestApplyChanges_AddAndRemove(t *testing.T) {
	out := applyChanges([]string{"prod"}, []string{"canary"}, []string{"prod"})
	if containsTag(out, "prod") {
		t.Fatalf("prod should be removed")
	}
	if !containsTag(out, "canary") {
		t.Fatalf("canary should be added")
	}
}

func TestTag_NilClientErrors(t *testing.T) {
	_, err := Tag(nil, "secret/data/app", []string{"prod"}, nil, true)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestTag_EmptyPathErrors(t *testing.T) {
	_, err := Tag(nil, "", []string{"prod"}, nil, true)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func containsTag(tags []string, target string) bool {
	for _, t := range tags {
		if t == target {
			return true
		}
	}
	return false
}
