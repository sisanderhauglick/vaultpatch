package group

import (
	"errors"
	"testing"
)

// stubClient implements VaultClient for tests.
type stubClient struct {
	data    map[string]map[string]interface{}
	written map[string]map[string]interface{}
	readErr error
}

func (s *stubClient) Read(path string) (map[string]interface{}, error) {
	if s.readErr != nil {
		return nil, s.readErr
	}
	return s.data[path], nil
}

func (s *stubClient) Write(path string, data map[string]interface{}) error {
	if s.written == nil {
		s.written = make(map[string]map[string]interface{})
	}
	s.written[path] = data
	return nil
}

func TestGroup_NilClientErrors(t *testing.T) {
	_, err := Group(nil, Options{Sources: []string{"a"}, Destination: "d"})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestGroup_EmptySourcesErrors(t *testing.T) {
	c := &stubClient{}
	_, err := Group(c, Options{Destination: "d"})
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestGroup_EmptyDestinationErrors(t *testing.T) {
	c := &stubClient{}
	_, err := Group(c, Options{Sources: []string{"a"}})
	if err == nil {
		t.Fatal("expected error for empty destination")
	}
}

func TestGroup_DryRun_ResultFlagged(t *testing.T) {
	c := &stubClient{
		data: map[string]map[string]interface{}{
			"secret/a": {"x": "1"},
		},
	}
	res, err := Group(c, Options{
		Sources:     []string{"secret/a"},
		Destination: "secret/out",
		DryRun:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.DryRun {
		t.Error("expected DryRun to be true")
	}
	if c.written["secret/out"] != nil {
		t.Error("expected no write in dry-run mode")
	}
}

func TestGroup_KeysMergedCount(t *testing.T) {
	c := &stubClient{
		data: map[string]map[string]interface{}{
			"secret/a": {"k1": "v1", "k2": "v2"},
			"secret/b": {"k3": "v3"},
		},
	}
	res, err := Group(c, Options{
		Sources:     []string{"secret/a", "secret/b"},
		Destination: "secret/out",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.KeysMerged != 3 {
		t.Errorf("expected 3 keys merged, got %d", res.KeysMerged)
	}
}

func TestGroup_ReadError_Propagates(t *testing.T) {
	c := &stubClient{readErr: errors.New("vault down")}
	_, err := Group(c, Options{
		Sources:     []string{"secret/a"},
		Destination: "secret/out",
	})
	if err == nil {
		t.Fatal("expected error from read failure")
	}
}

func TestPrefixKey_AppendsSegment(t *testing.T) {
	got := prefixKey("secret/env/prod", "db_pass")
	want := "prod/db_pass"
	if got != want {
		t.Errorf("prefixKey = %q, want %q", got, want)
	}
}
