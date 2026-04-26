package stamp_test

import (
	"github.com/youorg/vaultpatch/internal/vault"
)

// newStubClient returns a minimal *vault.Client suitable for unit tests.
// It relies on the exported constructor accepting explicit parameters so that
// no real Vault server is required.
func newStubClient() *vault.Client {
	c, err := vault.NewClient(vault.Params{
		Addr:  "http://127.0.0.1:8200",
		Token: "test-token",
	})
	if err != nil {
		panic("stamp_stub_test: failed to build stub client: " + err.Error())
	}
	return c
}
