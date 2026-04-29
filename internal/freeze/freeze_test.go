package freeze_test

import (
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"

	"github.com/your-org/vaultpatch/internal/freeze"
)

// stubClient satisfies freeze.Client for testing.
type stubClient struct {
	data    map[string]map[string]interface{}
	readErr error
	writeErr error
}

func (s *stubClient) Read(path string) (*api.Secret, error) {
	if s.readErr != nil {
		return nil, s.readErr
	}
	d, ok := s.data[path]
	if !ok {
		return nil, nil
	}
	return &api.Secret{Data: d}, nil
}

func (s *stubClient) Write(path string, data map[string]interface{}) (*api.Secret, error) {
	if s.writeErr != nil {
		return nil, s.writeErr
	}
	if s.data == nil {
		s.data = make(map[string]map[string]interface{})
	}
	s.data[path] = data
	return &api.Secret{Data: data}, nil
}

func TestFreeze_NilClientErrors(t *testing.T) {
	_, err := freeze.Freeze(nil, freeze.Options{Paths: []string{"secret/a"}, Owner: "ci"})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestFreeze_EmptyPathsError(t *testing.T) {
	c := &stubClient{}
	_, err := freeze.Freeze(c, freeze.Options{Owner: "ci"})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestFreeze_EmptyOwnerError(t *testing.T) {
	c := &stubClient{}
	_, err := freeze.Freeze(c, freeze.Options{Paths: []string{"secret/a"}})
	if err == nil {
		t.Fatal("expected error for empty owner")
	}
}

func TestFreeze_DryRun_ResultFlagged(t *testing.T) {
	c := &stubClient{}
	results, err := freeze.Freeze(c, freeze.Options{
		Paths:  []string{"secret/x"},
		Owner:  "alice",
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].DryRun {
		t.Error("expected DryRun=true")
	}
	if !results[0].Frozen {
		t.Error("expected Frozen=true")
	}
}

func TestFreeze_DryRun_NoWrite(t *testing.T) {
	c := &stubClient{}
	_, _ = freeze.Freeze(c, freeze.Options{
		Paths:  []string{"secret/y"},
		Owner:  "bob",
		DryRun: true,
	})
	if _, exists := c.data["secret/y"]; exists {
		t.Error("dry run should not write to vault")
	}
}

func TestUnfreeze_NilClientErrors(t *testing.T) {
	_, err := freeze.Unfreeze(nil, freeze.Options{Paths: []string{"secret/a"}})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestUnfreeze_EmptyPathsError(t *testing.T) {
	c := &stubClient{}
	_, err := freeze.Unfreeze(c, freeze.Options{})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestUnfreeze_DryRun_FrozenFalse(t *testing.T) {
	c := &stubClient{}
	results, err := freeze.Unfreeze(c, freeze.Options{
		Paths:  []string{"secret/z"},
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Frozen {
		t.Error("expected Frozen=false after unfreeze")
	}
}

func TestIsFrozen_ReturnsTrueWhenMarkerPresent(t *testing.T) {
	c := &stubClient{
		data: map[string]map[string]interface{}{
			"secret/locked": {"__frozen__": "alice@2024-01-01T00:00:00Z", "key": "val"},
		},
	}
	ok, err := freeze.IsFrozen(c, "secret/locked")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected IsFrozen=true")
	}
}

func TestIsFrozen_ReturnsFalseWhenMarkerAbsent(t *testing.T) {
	c := &stubClient{
		data: map[string]map[string]interface{}{
			"secret/open": {"key": "val"},
		},
	}
	ok, err := freeze.IsFrozen(c, "secret/open")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected IsFrozen=false")
	}
}

func TestFreeze_FrozenAtSet(t *testing.T) {
	c := &stubClient{}
	before := time.Now().UTC().Add(-time.Second)
	results, err := freeze.Freeze(c, freeze.Options{
		Paths:  []string{"secret/ts"},
		Owner:  "ci",
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].FrozenAt.Before(before) {
		t.Error("FrozenAt should be set to approximately now")
	}
}

func TestIsFrozen_ReadError(t *testing.T) {
	c := &stubClient{readErr: errors.New("vault down")}
	_, err := freeze.IsFrozen(c, "secret/any")
	if err == nil {
		t.Fatal("expected error on read failure")
	}
}
