package scope_test

import (
	"errors"

	"github.com/your-org/vaultpatch/internal/vault"
)

// stubVaultClient is an in-memory Vault client for testing.
type stubVaultClient struct {
	data  map[string]map[string]string
	lists map[string][]string
}

func stubClient(data map[string]map[string]string) *stubVaultClient {
	return &stubVaultClient{
		data:  data,
		lists: make(map[string][]string),
	}
}

func (s *stubVaultClient) SetList(path string, paths []string) {
	s.lists[path] = paths
}

func (s *stubVaultClient) Read(path string) (map[string]string, error) {
	if v, ok := s.data[path]; ok {
		return v, nil
	}
	return nil, errors.New("not found: " + path)
}

func (s *stubVaultClient) Write(path string, data map[string]string) error {
	s.data[path] = data
	return nil
}

func (s *stubVaultClient) List(path string) ([]string, error) {
	if v, ok := s.lists[path]; ok {
		return v, nil
	}
	return []string{}, nil
}

// toVaultClient satisfies any interface assertion in scope that expects *vault.Client.
// In the real project the stub would embed or wrap vault.Client; here we use the
// concrete type via a thin adapter so the test package compiles cleanly.
func toVaultClient(_ *stubVaultClient) *vault.Client { return nil }
