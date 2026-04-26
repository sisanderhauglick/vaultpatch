package stamp_test

import (
	"testing"

	"github.com/youorg/vaultpatch/internal/stamp"
)

func TestStamp_NilClientErrors(t *testing.T) {
	_, err := stamp.Stamp(nil, stamp.Options{
		Paths:       []string{"secret/app"},
		Annotations: map[string]string{"owner": "team-a"},
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestStamp_EmptyPathsError(t *testing.T) {
	client := newStubClient()
	_, err := stamp.Stamp(client, stamp.Options{
		Paths:       nil,
		Annotations: map[string]string{"env": "prod"},
	})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestStamp_EmptyAnnotationsError(t *testing.T) {
	client := newStubClient()
	_, err := stamp.Stamp(client, stamp.Options{
		Paths:       []string{"secret/app"},
		Annotations: nil,
	})
	if err == nil {
		t.Fatal("expected error for empty annotations")
	}
}

func TestStamp_DryRun_ResultFlagged(t *testing.T) {
	client := newStubClient()
	results, err := stamp.Stamp(client, stamp.Options{
		Paths:       []string{"secret/app"},
		Annotations: map[string]string{"owner": "team-a"},
		DryRun:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].DryRun {
		t.Error("expected DryRun to be true")
	}
	if results[0].Stamped {
		t.Error("expected Stamped to be false in dry-run")
	}
}

func TestStamp_DryRun_AnnotationsPreserved(t *testing.T) {
	client := newStubClient()
	annotations := map[string]string{"env": "staging", "team": "platform"}
	results, err := stamp.Stamp(client, stamp.Options{
		Paths:       []string{"secret/svc"},
		Annotations: annotations,
		DryRun:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, want := range annotations {
		if got := results[0].Annotations[k]; got != want {
			t.Errorf("annotation %q: got %q, want %q", k, got, want)
		}
	}
}

func TestStamp_DryRun_StampedAtSet(t *testing.T) {
	client := newStubClient()
	results, err := stamp.Stamp(client, stamp.Options{
		Paths:       []string{"secret/app"},
		Annotations: map[string]string{"owner": "ops"},
		DryRun:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].StampedAt.IsZero() {
		t.Error("expected StampedAt to be set")
	}
}
