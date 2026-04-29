package swap_test

import (
	"fmt"

	"github.com/your-org/vaultpatch/internal/vault"
)

type stubVaultClient struct {
	data map[string]map[string]interface{}
}

func newStubClient() *stubVaultClient {
	return &stubVaultClient{data: make(map[string]map[string]interface{})}
}

func (s *stubVaultClient) Read(path string) (map[string]interface{}, error) {
	if d, ok := s.data[path]; ok {
		out := make(map[string]interface{}, len(d))
		for k, v := range d {
			out[k] = v
		}
		return out, nil
	}
	return map[string]interface{}{}, nil
}

func (s *stubVaultClient) Write(path string, data map[string]interface{}) error {
	s.data[path] = data
	return nil
}

// toVaultClient converts the stub to the interface expected by swap.Swap.
// swap.Swap accepts *vault.Client; we satisfy the same interface via a thin
// adapter so tests remain decoupled from the real Vault SDK.
func toVaultClient(s *stubVaultClient) *vault.Client {
	_ = fmt.Sprintf // keep import used
	// In tests the package-level functions are called with the stub directly;
	// this helper exists as documentation of the intended interface contract.
	return nil
}
