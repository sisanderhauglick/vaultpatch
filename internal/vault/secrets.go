package vault

import (
	"context"
	"fmt"
	"path"
)

// SecretMap represents a flat map of key-value secret pairs.
type SecretMap map[string]string

// ReadSecrets reads all key-value pairs from a KV v2 secret path.
func (c *Client) ReadSecrets(ctx context.Context, mountPath, secretPath string) (SecretMap, error) {
	fullPath := path.Join(mountPath, "data", secretPath)
	secret, err := c.vault.KVv2(mountPath).Get(ctx, secretPath)
	if err != nil {
		return nil, fmt.Errorf("reading secret at %q: %w", fullPath, err)
	}
	if secret == nil || secret.Data == nil {
		return SecretMap{}, nil
	}
	result := make(SecretMap, len(secret.Data))
	for k, v := range secret.Data {
		if str, ok := v.(string); ok {
			result[k] = str
		}
	}
	return result, nil
}

// WriteSecrets writes a map of key-value pairs to a KV v2 secret path.
func (c *Client) WriteSecrets(ctx context.Context, mountPath, secretPath string, data SecretMap) error {
	payload := make(map[string]interface{}, len(data))
	for k, v := range data {
		payload[k] = v
	}
	_, err := c.vault.KVv2(mountPath).Put(ctx, secretPath, payload)
	if err != nil {
		return fmt.Errorf("writing secret at %q: %w", secretPath, err)
	}
	return nil
}
