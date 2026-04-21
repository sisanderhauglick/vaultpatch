package sanitize_test

import (
	"github.com/your-org/vaultpatch/internal/vault"
)

// stubClient creates a minimal *vault.Client backed by in-memory data for
// testing purposes. It returns the client and the underlying secret map so
// callers can inspect writes.
func stubClient(data map[string]interface{}) (*vault.Client, map[string]interface{}) {
	copy := make(map[string]interface{}, len(data))
	for k, v := range data {
		copy[k] = v
	}
	return vault.NewStubClient(copy), copy
}
