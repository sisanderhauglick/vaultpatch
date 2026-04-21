package flatten_test

import (
	"testing"

	"github.com/yourusername/vaultpatch/internal/flatten"
)

func TestFlatten_NilClientErrors(t *testing.T) {
	_, err := flatten.Flatten(nil, flatten.Options{
		Sources:     []string{"secret/a"},
		Destination: "secret/dest",
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestFlatten_EmptySourcesErrors(t *testing.T) {
	_, err := flatten.Flatten(stubClient(), flatten.Options{
		Sources:     []string{},
		Destination: "secret/dest",
	})
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestFlatten_EmptyDestinationErrors(t *testing.T) {
	_, err := flatten.Flatten(stubClient(), flatten.Options{
		Sources:     []string{"secret/a"},
		Destination: "",
	})
	if err == nil {
		t.Fatal("expected error for empty destination")
	}
}

func TestFlatten_DryRun_ResultFlagged(t *testing.T) {
	res, err := flatten.Flatten(stubClient(), flatten.Options{
		Sources:     []string{"secret/a"},
		Destination: "secret/dest",
		DryRun:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestFlatten_DryRun_KeysMergedCount(t *testing.T) {
	res, err := flatten.Flatten(stubClient(), flatten.Options{
		Sources:     []string{"secret/a", "secret/b"},
		Destination: "secret/dest",
		DryRun:      true,
		Overwrite:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// stubClient provides 2 keys per path; with overwrite, all 4 should be present
	if res.KeysMerged != 4 {
		t.Errorf("expected 4 keys merged, got %d", res.KeysMerged)
	}
}

func TestFlatten_DestinationPreserved(t *testing.T) {
	res, err := flatten.Flatten(stubClient(), flatten.Options{
		Sources:     []string{"secret/a"},
		Destination: "secret/dest",
		DryRun:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Destination != "secret/dest" {
		t.Errorf("expected destination %q, got %q", "secret/dest", res.Destination)
	}
}

func TestFlatten_FlattenedAtSet(t *testing.T) {
	res, err := flatten.Flatten(stubClient(), flatten.Options{
		Sources:     []string{"secret/a"},
		Destination: "secret/dest",
		DryRun:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.FlattenedAt.IsZero() {
		t.Error("expected FlattenedAt to be set")
	}
}
