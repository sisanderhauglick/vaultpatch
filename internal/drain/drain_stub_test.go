package drain_test

import (
	"fmt"

	"github.com/your-org/vaultpatch/internal/vault"
)

// stubStore is an in-memory Vault backend used by tests.
type stubStore struct {
	data map[string]map[string]interface{}
}

func newStubClient(data map[string]map[string]interface{}) *vault.Client {
	s := &stubStore{data: data}
	return vault.NewClientFromRW(s.read, s.write)
}

func (s *stubStore) read(path string) (map[string]interface{}, error) {
	v, ok := s.data[path]
	if !ok {
		return map[string]interface{}{}, nil
	}
	// Return a shallow copy so mutations are visible only after write.
	copy := make(map[string]interface{}, len(v))
	for k, val := range v {
		copy[k] = val
	}
	return copy, nil
}

func (s *stubStore) write(path string, data map[string]interface{}) error {
	if s.data == nil {
		return fmt.Errorf("stub: no store initialised")
	}
	next := make(map[string]interface{}, len(data))
	for k, v := range data {
		next[k] = v
	}
	s.data[path] = next
	return nil
}
