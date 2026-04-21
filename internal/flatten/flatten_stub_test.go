package flatten_test

import (
	"fmt"

	"github.com/yourusername/vaultpatch/internal/vault"
)

// stubData maps path -> key/value pairs used by the stub client.
var stubData = map[string]map[string]interface{}{
	"secret/a": {
		"key_a1": "val_a1",
		"key_a2": "val_a2",
	},
	"secret/b": {
		"key_b1": "val_b1",
		"key_b2": "val_b2",
	},
	"secret/dest": {},
}

// stubClient returns a *vault.Client configured to talk to a fake server that
// serves stubData. In tests the vault package's ReadSecrets / WriteSecrets are
// exercised through the real code paths; we rely on the fake Vault server
// pattern already established in internal/vault/client_test.go.
//
// For unit-testing purposes we use a minimal httptest-backed stub so the
// flatten logic can be exercised without a real Vault instance.
func stubClient() *vault.Client {
	c, err := vault.NewClient(vault.Params{
		Addr:  fakeVaultServer().URL,
		Token: "test-token",
	})
	if err != nil {
		panic(fmt.Sprintf("stubClient: %v", err))
	}
	return c
}
