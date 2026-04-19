package pin_test

import "github.com/your-org/vaultpatch/internal/vault"

// stubClient satisfies vault.Client's Read/Write interface for unit tests.
type stubClient struct{}

func (s *stubClient) Read(_ string) (map[string]string, error) {
	return map[string]string{"key": "value"}, nil
}

func (s *stubClient) Write(_ string, _ map[string]string) error {
	return nil
}

// Ensure stubClient matches the shape expected by pin package.
var _ interface {
	Read(string) (map[string]string, error)
	Write(string, map[string]string) error
} = (*stubClient)(nil)

// vaultClientShim adapts stubClient to *vault.Client for tests that need it.
func toVaultClient(_ *stubClient) *vault.Client {
	return nil
}
