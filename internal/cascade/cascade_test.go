package cascade

import (
	"errors"
	"testing"

	"github.com/hashicorp/vault/api"
)

type stubClient struct {
	data    map[string]interface{}
	readErr error
	writeErr error
	written map[string]map[string]interface{}
}

func (s *stubClient) Read(_ string) (*api.Secret, error) {
	if s.readErr != nil {
		return nil, s.readErr
	}
	return &api.Secret{Data: s.data}, nil
}

func (s *stubClient) Write(path string, data map[string]interface{}) (*api.Secret, error) {
	if s.writeErr != nil {
		return nil, s.writeErr
	}
	if s.written == nil {
		s.written = make(map[string]map[string]interface{})
	}
	s.written[path] = data
	return &api.Secret{}, nil
}

func TestCascade_NilClientErrors(t *testing.T) {
	_, err := Cascade(nil, Options{Source: "src", Destinations: []string{"dst"}})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestCascade_EmptySourceErrors(t *testing.T) {
	_, err := Cascade(&stubClient{}, Options{Destinations: []string{"dst"}})
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestCascade_EmptyDestinationsErrors(t *testing.T) {
	_, err := Cascade(&stubClient{}, Options{Source: "src"})
	if err == nil {
		t.Fatal("expected error for empty destinations")
	}
}

func TestCascade_DryRun_SkipsWrite(t *testing.T) {
	c := &stubClient{data: map[string]interface{}{"k": "v"}}
	results, err := Cascade(c, Options{
		Source:       "secret/src",
		Destinations: []string{"secret/dst1", "secret/dst2"},
		DryRun:       true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if len(c.written) != 0 {
		t.Error("expected no writes in dry-run mode")
	}
	for _, r := range results {
		if !r.DryRun {
			t.Error("expected DryRun flag to be true")
		}
	}
}

func TestCascade_LiveRun_WritesAllDestinations(t *testing.T) {
	c := &stubClient{data: map[string]interface{}{"a": "1", "b": "2"}}
	dests := []string{"secret/d1", "secret/d2", "secret/d3"}
	results, err := Cascade(c, Options{
		Source:       "secret/src",
		Destinations: dests,
		DryRun:       false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, d := range dests {
		if _, ok := c.written[d]; !ok {
			t.Errorf("expected write to %q", d)
		}
	}
}

func TestCascade_KeyFilter_SubsetOnly(t *testing.T) {
	c := &stubClient{data: map[string]interface{}{"x": "1", "y": "2", "z": "3"}}
	_, err := Cascade(c, Options{
		Source:       "secret/src",
		Destinations: []string{"secret/dst"},
		Keys:         []string{"x", "z"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	written := c.written["secret/dst"]
	if len(written) != 2 {
		t.Fatalf("expected 2 keys written, got %d", len(written))
	}
	if _, ok := written["y"]; ok {
		t.Error("key 'y' should not have been cascaded")
	}
}

func TestCascade_ReadError_Propagates(t *testing.T) {
	c := &stubClient{readErr: errors.New("vault unavailable")}
	_, err := Cascade(c, Options{
		Source:       "secret/src",
		Destinations: []string{"secret/dst"},
	})
	if err == nil {
		t.Fatal("expected read error to propagate")
	}
}

func TestResult_CascadedAtSet(t *testing.T) {
	c := &stubClient{data: map[string]interface{}{"k": "v"}}
	results, err := Cascade(c, Options{
		Source:       "secret/src",
		Destinations: []string{"secret/dst"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].CascadedAt.IsZero() {
		t.Error("expected CascadedAt to be set")
	}
}
