package pin_test

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/pin"
)

func TestPin_NilClientErrors(t *testing.T) {
	_, err := pin.Pin(nil, pin.Options{Paths: []string{"secret/a"}, Version: 1})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestPin_EmptyPathsError(t *testing.T) {
	_, err := pin.Pin(&stubClient{}, pin.Options{Version: 1})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestPin_InvalidVersionError(t *testing.T) {
	_, err := pin.Pin(&stubClient{}, pin.Options{Paths: []string{"secret/a"}, Version: 0})
	if err == nil {
		t.Fatal("expected error for version < 1")
	}
}

func TestPin_DryRun_ResultFlagged(t *testing.T) {
	results, err := pin.Pin(&stubClient{}, pin.Options{
		Paths:   []string{"secret/a", "secret/b"},
		Version: 3,
		DryRun:  true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.DryRun {
			t.Errorf("expected DryRun=true for %s", r.Path)
		}
		if !r.Pinned {
			t.Errorf("expected Pinned=true for %s", r.Path)
		}
		if r.Version != 3 {
			t.Errorf("expected Version=3 for %s, got %d", r.Path, r.Version)
		}
	}
}

func TestUnpin_NilClientErrors(t *testing.T) {
	_, err := pin.Unpin(nil, []string{"secret/a"}, true)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestUnpin_EmptyPathsError(t *testing.T) {
	_, err := pin.Unpin(&stubClient{}, []string{}, false)
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestUnpin_DryRun_NotPinned(t *testing.T) {
	results, err := pin.Unpin(&stubClient{}, []string{"secret/a"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if results[0].Pinned {
		t.Error("expected Pinned=false after unpin")
	}
	if !results[0].DryRun {
		t.Error("expected DryRun=true")
	}
}
