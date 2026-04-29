package blend_test

import (
	"errors"
	"testing"

	"github.com/hashicorp/vault/api"

	"github.com/your-org/vaultpatch/internal/blend"
)

type stubClient struct {
	data    map[string]map[string]interface{}
	written map[string]map[string]interface{}
	readErr error
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
	if s.written == nil {
		s.written = make(map[string]map[string]interface{})
	}
	s.written[path] = data
	return &api.Secret{}, nil
}

func TestBlend_NilClientErrors(t *testing.T) {
	_, err := blend.Blend(nil, blend.Options{Sources: []string{"a"}, Destination: "dest"})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestBlend_EmptySourcesErrors(t *testing.T) {
	c := &stubClient{}
	_, err := blend.Blend(c, blend.Options{Destination: "dest"})
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestBlend_EmptyDestinationErrors(t *testing.T) {
	c := &stubClient{}
	_, err := blend.Blend(c, blend.Options{Sources: []string{"a"}})
	if err == nil {
		t.Fatal("expected error for empty destination")
	}
}

func TestBlend_StrategyLast_OverwritesConflict(t *testing.T) {
	c := &stubClient{
		data: map[string]map[string]interface{}{
			"src/a": {"key": "first"},
			"src/b": {"key": "last"},
		},
	}
	res, err := blend.Blend(c, blend.Options{
		Sources:     []string{"src/a", "src/b"},
		Destination: "dst",
		Strategy:    blend.StrategyLast,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.written["dst"]["key"] != "last" {
		t.Errorf("expected 'last', got %v", c.written["dst"]["key"])
	}
	if res.Destination != "dst" {
		t.Errorf("unexpected destination: %s", res.Destination)
	}
}

func TestBlend_StrategyFirst_KeepsFirstConflict(t *testing.T) {
	c := &stubClient{
		data: map[string]map[string]interface{}{
			"src/a": {"key": "first"},
			"src/b": {"key": "last"},
		},
	}
	_, err := blend.Blend(c, blend.Options{
		Sources:     []string{"src/a", "src/b"},
		Destination: "dst",
		Strategy:    blend.StrategyFirst,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.written["dst"]["key"] != "first" {
		t.Errorf("expected 'first', got %v", c.written["dst"]["key"])
	}
}

func TestBlend_DryRun_ResultFlagged(t *testing.T) {
	c := &stubClient{
		data: map[string]map[string]interface{}{
			"src/a": {"x": "1"},
		},
	}
	res, err := blend.Blend(c, blend.Options{
		Sources:     []string{"src/a"},
		Destination: "dst",
		DryRun:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.DryRun {
		t.Error("expected DryRun flag to be true")
	}
	if len(c.written) != 0 {
		t.Error("expected no writes during dry run")
	}
}

func TestBlend_ReadError_Propagates(t *testing.T) {
	c := &stubClient{readErr: errors.New("vault unavailable")}
	_, err := blend.Blend(c, blend.Options{
		Sources:     []string{"src/a"},
		Destination: "dst",
	})
	if err == nil {
		t.Fatal("expected error from read failure")
	}
}

func TestBlend_KeyFilter_OnlyIncludesSpecified(t *testing.T) {
	c := &stubClient{
		data: map[string]map[string]interface{}{
			"src/a": {"keep": "yes", "drop": "no"},
		},
	}
	_, err := blend.Blend(c, blend.Options{
		Sources:     []string{"src/a"},
		Destination: "dst",
		Keys:        []string{"keep"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := c.written["dst"]["drop"]; ok {
		t.Error("expected 'drop' key to be excluded")
	}
	if c.written["dst"]["keep"] != "yes" {
		t.Errorf("expected 'keep' key to be present")
	}
}
