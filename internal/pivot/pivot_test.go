package pivot

import (
	"errors"
	"sort"
	"testing"
)

// stubClient is a minimal in-memory VaultClient for tests.
type stubClient struct {
	store  map[string]map[string]interface{}
	writen map[string]map[string]interface{}
	readErr error
}

func (s *stubClient) Read(path string) (map[string]interface{}, error) {
	if s.readErr != nil {
		return nil, s.readErr
	}
	return s.store[path], nil
}

func (s *stubClient) Write(path string, data map[string]interface{}) error {
	if s.writen == nil {
		s.writen = make(map[string]map[string]interface{})
	}
	s.writen[path] = data
	return nil
}

func TestPivot_NilClientErrors(t *testing.T) {
	_, err := Pivot(nil, Options{Sources: []string{"a"}, Destination: "dst"})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestPivot_EmptySourcesErrors(t *testing.T) {
	_, err := Pivot(&stubClient{}, Options{Destination: "dst"})
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestPivot_EmptyDestinationErrors(t *testing.T) {
	_, err := Pivot(&stubClient{}, Options{Sources: []string{"a"}})
	if err == nil {
		t.Fatal("expected error for empty destination")
	}
}

func TestPivot_DryRun_ResultFlagged(t *testing.T) {
	client := &stubClient{
		store: map[string]map[string]interface{}{
			"secret/app": {"key": "val"},
		},
	}
	res, err := Pivot(client, Options{
		Sources:     []string{"secret/app"},
		Destination: "secret/merged",
		DryRun:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.DryRun {
		t.Error("expected DryRun to be true")
	}
	if len(client.writen) != 0 {
		t.Error("expected no writes during dry run")
	}
}

func TestPivot_KeysNamespacedByPathSegment(t *testing.T) {
	client := &stubClient{
		store: map[string]map[string]interface{}{
			"secret/svc/alpha": {"host": "a", "port": "1"},
			"secret/svc/beta":  {"host": "b"},
		},
	}
	res, err := Pivot(client, Options{
		Sources:     []string{"secret/svc/alpha", "secret/svc/beta"},
		Destination: "secret/combined",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sort.Strings(res.MergedKeys)
	want := []string{"alpha_host", "alpha_port", "beta_host"}
	for i, k := range want {
		if res.MergedKeys[i] != k {
			t.Errorf("key[%d]: got %q, want %q", i, res.MergedKeys[i], k)
		}
	}
}

func TestPivot_ReadError_Propagates(t *testing.T) {
	client := &stubClient{readErr: errors.New("vault unavailable")}
	_, err := Pivot(client, Options{
		Sources:     []string{"secret/x"},
		Destination: "secret/dst",
	})
	if err == nil {
		t.Fatal("expected error from read failure")
	}
}
