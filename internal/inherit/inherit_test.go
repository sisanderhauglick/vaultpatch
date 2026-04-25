package inherit

import (
	"errors"
	"testing"
)

type stubClient struct {
	store    map[string]map[string]string
	writeErr error
	readErr  error
}

func (s *stubClient) Read(path string) (map[string]string, error) {
	if s.readErr != nil {
		return nil, s.readErr
	}
	return s.store[path], nil
}

func (s *stubClient) Write(path string, data map[string]string) error {
	if s.writeErr != nil {
		return s.writeErr
	}
	if s.store == nil {
		s.store = map[string]map[string]string{}
	}
	s.store[path] = data
	return nil
}

func TestInherit_NilClientErrors(t *testing.T) {
	_, err := Inherit(nil, Options{Parent: "secret/parent", Children: []string{"secret/child"}})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestInherit_EmptyParentErrors(t *testing.T) {
	c := &stubClient{store: map[string]map[string]string{}}
	_, err := Inherit(c, Options{Children: []string{"secret/child"}})
	if err == nil {
		t.Fatal("expected error for empty parent")
	}
}

func TestInherit_EmptyChildrenErrors(t *testing.T) {
	c := &stubClient{store: map[string]map[string]string{}}
	_, err := Inherit(c, Options{Parent: "secret/parent"})
	if err == nil {
		t.Fatal("expected error for empty children")
	}
}

func TestInherit_DryRun_NoWrite(t *testing.T) {
	c := &stubClient{
		store: map[string]map[string]string{
			"secret/parent": {"DB_HOST": "prod-db", "DB_PORT": "5432"},
			"secret/child":  {},
		},
	}
	results, err := Inherit(c, Options{
		Parent:   "secret/parent",
		Children: []string{"secret/child"},
		DryRun:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].DryRun {
		t.Error("expected DryRun flag to be true")
	}
	// child store should remain empty
	if got := c.store["secret/child"]; len(got) != 0 {
		t.Errorf("expected no write in dry-run, got %v", got)
	}
}

func TestInherit_PropagatesKeys(t *testing.T) {
	c := &stubClient{
		store: map[string]map[string]string{
			"secret/parent": {"API_KEY": "xyz", "TIMEOUT": "30s"},
			"secret/child":  {},
		},
	}
	results, err := Inherit(c, Options{
		Parent:   "secret/parent",
		Children: []string{"secret/child"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results[0].Inherited) != 2 {
		t.Errorf("expected 2 inherited keys, got %d", len(results[0].Inherited))
	}
}

func TestInherit_SkipsExistingWithoutForce(t *testing.T) {
	c := &stubClient{
		store: map[string]map[string]string{
			"secret/parent": {"API_KEY": "parent-val"},
			"secret/child":  {"API_KEY": "child-val"},
		},
	}
	results, err := Inherit(c, Options{
		Parent:   "secret/parent",
		Children: []string{"secret/child"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results[0].Skipped) != 1 {
		t.Errorf("expected 1 skipped key, got %d", len(results[0].Skipped))
	}
}

func TestInherit_ForceOverwritesExisting(t *testing.T) {
	c := &stubClient{
		store: map[string]map[string]string{
			"secret/parent": {"API_KEY": "new-val"},
			"secret/child":  {"API_KEY": "old-val"},
		},
	}
	_, err := Inherit(c, Options{
		Parent:   "secret/parent",
		Children: []string{"secret/child"},
		Force:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := c.store["secret/child"]["API_KEY"]; got != "new-val" {
		t.Errorf("expected overwritten value %q, got %q", "new-val", got)
	}
}

func TestInherit_ReadError(t *testing.T) {
	c := &stubClient{readErr: errors.New("vault unavailable")}
	_, err := Inherit(c, Options{
		Parent:   "secret/parent",
		Children: []string{"secret/child"},
	})
	if err == nil {
		t.Fatal("expected error on read failure")
	}
}

func TestFilterKeys_SubsetOnly(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	out := filterKeys(secrets, []string{"A", "C"})
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["B"]; ok {
		t.Error("key B should have been excluded")
	}
}
