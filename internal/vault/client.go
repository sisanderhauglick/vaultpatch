package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the HashiCorp Vault API client.
type Client struct {
	v *vaultapi.Client
}

// NewClient creates a new Vault client using the provided address and token.
// If addr or token are empty, it falls back to VAULT_ADDR and VAULT_TOKEN env vars.
func NewClient(addr, token string) (*Client, error) {
	if addr == "" {
		addr = os.Getenv("VAULT_ADDR")
	}
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if addr == "" {
		return nil, fmt.Errorf("vault address is required (set VAULT_ADDR or pass --addr)")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required (set VAULT_TOKEN or pass --token)")
	}

	cfg := vaultapi.DefaultConfig()
	cfg.Address = addr

	v, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}
	v.SetToken(token)

	return &Client{v: v}, nil
}

// ReadSecret reads a KV v2 secret at the given mount and path.
// Returns a map of key/value pairs or an error.
func (c *Client) ReadSecret(mount, path string) (map[string]string, error) {
	fullPath := fmt.Sprintf("%s/data/%s", mount, path)
	secret, err := c.v.Logical().Read(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %q: %w", fullPath, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no secret found at %q", fullPath)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format at %q", fullPath)
	}

	result := make(map[string]string, len(data))
	for k, v := range data {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}

// WriteSecret writes key/value pairs to a KV v2 secret at the given mount and path.
func (c *Client) WriteSecret(mount, path string, data map[string]string) error {
	fullPath := fmt.Sprintf("%s/data/%s", mount, path)

	payload := make(map[string]interface{}, len(data))
	for k, v := range data {
		payload[k] = v
	}

	_, err := c.v.Logical().Write(fullPath, map[string]interface{}{"data": payload})
	if err != nil {
		return fmt.Errorf("failed to write secret at %q: %w", fullPath, err)
	}
	return nil
}
