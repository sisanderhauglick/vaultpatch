package prune_test

import (
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultpatch/internal/prune"
)

// stubClient implements prune.Client for tests.
type stubClient struct {
	data    map[string]map[string]interface{}
	deleted []string
	readErr error
}

func (s *stubClient) Read(path string) (map[string]interface{}, error) {
	if s.readErr != nil {
		return nil, s.readErr
	}
	return s.data[path], nil
}

func (s *stubClient) Write(path string, data map[string]interface{}) error { return nil }

func (s *stubClient) Delete(path string) error {
	s.deleted = append(s.deleted, path)
	return nil
}

func oldTimestamp() string {
	return time.Now().UTC().Add(-48 * time.Hour).Format(time.RFC3339)
}

func recentTimestamp() string {
	return time.Now().UTC().Add(-1 * time.Hour).Format(time.RFC3339)
}

func TestPrune_NilClientErrors(t *testing.T) {
	_, err := prune.Prune(nil, prune.Options{Paths: []string{"secret/a"}, OlderThan: time.Hour})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestPrune_EmptyPathsError(t *testing.T) {
	c := &stubClient{}
	_, err := prune.Prune(c, prune.Options{OlderThan: time.Hour})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestPrune_ZeroDurationError(t *testing.T) {
	c := &stubClient{}
	_, err := prune.Prune(c, prune.Options{Paths: []string{"secret/a"}})
	if err == nil {
		t.Fatal("expected error for zero OlderThan")
	}
}

func TestPrune_OldSecret_Pruned(t *testing.T) {
	c := &stubClient{
		data: map[string]map[string]interface{}{
			"secret/old": {"_created_at": oldTimestamp(), "key": "val"},
		},
	}
	results, err := prune.Prune(c, prune.Options{
		Paths:     []string{"secret/old"},
		OlderThan: 24 * time.Hour,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Pruned {
		t.Errorf("expected path to be pruned")
	}
	if len(c.deleted) != 1 {
		t.Errorf("expected Delete to be called once, got %d", len(c.deleted))
	}
}

func TestPrune_RecentSecret_Skipped(t *testing.T) {
	c := &stubClient{
		data: map[string]map[string]interface{}{
			"secret/new": {"_created_at": recentTimestamp()},
		},
	}
	results, err := prune.Prune(c, prune.Options{
		Paths:     []string{"secret/new"},
		OlderThan: 24 * time.Hour,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Pruned {
		t.Errorf("expected recent secret to be skipped")
	}
}

func TestPrune_DryRun_DoesNotDelete(t *testing.T) {
	c := &stubClient{
		data: map[string]map[string]interface{}{
			"secret/old": {"_created_at": oldTimestamp()},
		},
	}
	results, err := prune.Prune(c, prune.Options{
		Paths:     []string{"secret/old"},
		OlderThan: 24 * time.Hour,
		DryRun:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Pruned {
		t.Errorf("expected DryRun result to report Pruned=true")
	}
	if !results[0].DryRun {
		t.Errorf("expected DryRun flag set")
	}
	if len(c.deleted) != 0 {
		t.Errorf("expected no actual deletes during dry run")
	}
}

func TestPrune_MissingCreatedAt_Skipped(t *testing.T) {
	c := &stubClient{
		data: map[string]map[string]interface{}{
			"secret/nodates": {"key": "value"},
		},
	}
	results, err := prune.Prune(c, prune.Options{
		Paths:     []string{"secret/nodates"},
		OlderThan: time.Hour,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Pruned {
		t.Errorf("expected path without _created_at to be skipped")
	}
}

func TestPrune_ReadError_Propagated(t *testing.T) {
	c := &stubClient{readErr: errors.New("vault unavailable")}
	_, err := prune.Prune(c, prune.Options{
		Paths:     []string{"secret/x"},
		OlderThan: time.Hour,
	})
	if err == nil {
		t.Fatal("expected read error to propagate")
	}
}
