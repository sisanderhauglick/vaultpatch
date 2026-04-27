package prefill

import (
	"github.com/youorg/vaultpatch/internal/vault"
)

// stubClient returns a *vault.Client wired to an in-memory store with no
// pre-existing secrets.
func stubClient() *vault.Client {
	return stubClientWithData(nil)
}

// stubClientWithData returns a *vault.Client whose read calls return data.
func stubClientWithData(data map[string]interface{}) *vault.Client {
	if data == nil {
		data = map[string]interface{}{}
	}
	return vault.NewClientFromStore(
		vault.NewMemStore(data),
	)
}
