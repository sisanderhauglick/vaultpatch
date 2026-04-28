package enrich_test

import (
	"errors"
	"testing"

	"github.com/vaultpatch/vaultpatch/internal/enrich"
)

// stubClient is a minimal in-memory VaultClient.
type stubClient struct {
	store   map[string]map[string]interface{}
	readErr error
	writeErr error
}

func (s *stubClient) Read(path string) (map[string]interface{}, error) {
	if s.readErr != nil {
		return nil, s.readErr
	}
	return s.store[path], nil
}

func (s *stubClient) Write(path string, data map[string]interface{}) error {
	if s.writeErr != nil {
		return s.writeErr
	}
	if s.store == nil {
		s.store = make(map[string]map[string]interface{})
	}
	s.store[path] = data
	return nil
}

func TestEnrich_NilClientErrors(t *testing.T) {
	_, err := enrich.Enrich(nil, enrich.Options{Paths: []string{"secret/a"}, Annotations: map[string]string{"k": "v"}})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestEnrich_EmptyPathsError(t *testing.T) {
	c := &stubClient{}
	_, err := enrich.Enrich(c, enrich.Options{Annotations: map[string]string{"k": "v"}})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestEnrich_EmptyAnnotationsError(t *testing.T) {
	c := &stubClient{}
	_, err := enrich.Enrich(c, enrich.Options{Paths: []string{"secret/a"}})
	if err == nil {
		t.Fatal("expected error for empty annotations")
	}
}

func TestEnrich_DryRun_ResultFlagged(t *testing.T) {
	c := &stubClient{store: map[string]map[string]interface{}{"secret/a": {"existing": "val"}}}
	res, err := enrich.Enrich(c, enrich.Options{
		Paths:       []string{"secret/a"},
		Annotations: map[string]string{"env": "prod"},
		DryRun:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 || !res[0].DryRun {
		t.Fatal("expected DryRun=true in result")
	}
	// store should be unchanged on dry-run
	if _, ok := c.store["secret/a"]["env"]; ok {
		t.Fatal("dry-run must not write to store")
	}
}

func TestEnrich_LiveRun_WritesAnnotations(t *testing.T) {
	c := &stubClient{store: map[string]map[string]interface{}{"secret/b": {"foo": "bar"}}}
	res, err := enrich.Enrich(c, enrich.Options{
		Paths:       []string{"secret/b"},
		Annotations: map[string]string{"owner": "team-a"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) == 0 || res[0].DryRun {
		t.Fatal("expected live result")
	}
	if c.store["secret/b"]["owner"] != "team-a" {
		t.Fatal("expected annotation to be written")
	}
}

func TestEnrich_ReadError_Propagates(t *testing.T) {
	c := &stubClient{readErr: errors.New("vault unavailable")}
	_, err := enrich.Enrich(c, enrich.Options{
		Paths:       []string{"secret/c"},
		Annotations: map[string]string{"k": "v"},
	})
	if err == nil {
		t.Fatal("expected read error to propagate")
	}
}
