package immutable_test

import (
	"errors"
	"testing"

	"github.com/yourusername/vaultpatch/internal/immutable"
)

// stubClient satisfies VaultClient for tests.
type stubClient struct {
	store   map[string]map[string]interface{}
	readErr error
	writeErr error
}

func newStub() *stubClient {
	return &stubClient{store: make(map[string]map[string]interface{})}
}

func (s *stubClient) Read(path string) (map[string]interface{}, error) {
	if s.readErr != nil {
		return nil, s.readErr
	}
	if v, ok := s.store[path]; ok {
		return v, nil
	}
	return map[string]interface{}{}, nil
}

func (s *stubClient) Write(path string, data map[string]interface{}) error {
	if s.writeErr != nil {
		return s.writeErr
	}
	s.store[path] = data
	return nil
}

func TestImmute_NilClientErrors(t *testing.T) {
	_, err := immutable.Immute(nil, immutable.Options{Paths: []string{"secret/a"}})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestImmute_EmptyPathsError(t *testing.T) {
	_, err := immutable.Immute(newStub(), immutable.Options{})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestImmute_DryRun_ResultFlagged(t *testing.T) {
	c := newStub()
	results, err := immutable.Immute(c, immutable.Options{Paths: []string{"secret/a"}, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].DryRun {
		t.Error("expected DryRun to be true")
	}
	if _, ok := c.store["secret/a"]; ok {
		t.Error("dry run should not write to store")
	}
}

func TestImmute_LiveRun_Writessentinel(t *testing.T) {
	c := newStub()
	_, err := immutable.Immute(c, immutable.Options{Paths: []string{"secret/b"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, ok := c.store["secret/b"]
	if !ok {
		t.Fatal("expected data to be written")
	}
	if _, found := data["__immutable"]; !found {
		t.Error("expected __immutable key to be present")
	}
}

func TestRelease_NilClientErrors(t *testing.T) {
	_, err := immutable.Release(nil, immutable.Options{Paths: []string{"secret/a"}})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestRelease_DryRun_DoesNotWrite(t *testing.T) {
	c := newStub()
	c.store["secret/c"] = map[string]interface{}{"__immutable": "2024-01-01T00:00:00Z", "key": "val"}
	_, err := immutable.Release(c, immutable.Options{Paths: []string{"secret/c"}, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, still := c.store["secret/c"]["__immutable"]; !still {
		t.Error("dry run should not remove sentinel")
	}
}

func TestIsImmutable_NilClientErrors(t *testing.T) {
	_, err := immutable.IsImmutable(nil, "secret/x")
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestIsImmutable_DetectsSentinel(t *testing.T) {
	c := newStub()
	c.store["secret/d"] = map[string]interface{}{"__immutable": "2024-01-01T00:00:00Z"}
	ok, err := immutable.IsImmutable(c, "secret/d")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected path to be detected as immutable")
	}
}

func TestIsImmutable_ReadError(t *testing.T) {
	c := newStub()
	c.readErr = errors.New("vault unavailable")
	_, err := immutable.IsImmutable(c, "secret/e")
	if err == nil {
		t.Fatal("expected error on read failure")
	}
}
