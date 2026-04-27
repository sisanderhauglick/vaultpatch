package reorder

import (
	"github.com/vaultpatch/vaultpatch/internal/vault"
)

// stubClient returns a minimal *vault.Client suitable for unit tests that do
// not exercise real Vault connectivity. The address and token values satisfy
// the client constructor validation without requiring a live server.
func stubClient() *vault.Client {
	c, err := vault.NewClient(vault.Params{
		Addr:  "http://127.0.0.1:8200",
		Token: "test-token",
	})
	if err != nil {
		panic("reorder stub: failed to create client: " + err.Error())
	}
	return c
}
