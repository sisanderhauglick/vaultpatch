// Package rollback provides functionality to revert Vault secrets to a previous state.
package rollback

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Snapshot holds a point-in-time copy of secrets at a given path.
type Snapshot struct {
	Path    string
	Secrets map[string]string
}

// Capture reads the current secrets at path and returns a Snapshot.
func Capture(ctx context.Context, client *api.Client, path string) (*Snapshot, error) {
	if client == nil {
		return nil, errors.New("rollback: vault client is nil")
	}
	if path == "" {
		return nil, errors.New("rollback: path must not be empty")
	}

	secret, err := client.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("rollback: read %s: %w", path, err)
	}

	data := make(map[string]string)
	if secret != nil && secret.Data != nil {
		for k, v := range secret.Data {
			if s, ok := v.(string); ok {
				data[k] = s
			}
		}
	}

	return &Snapshot{Path: path, Secrets: data}, nil
}

// Restore writes the snapshot secrets back to Vault.
// If dryRun is true, no writes are performed.
func Restore(ctx context.Context, client *api.Client, snap *Snapshot, dryRun bool) error {
	if client == nil {
		return errors.New("rollback: vault client is nil")
	}
	if snap == nil {
		return errors.New("rollback: snapshot is nil")
	}

	if dryRun {
		fmt.Printf("[dry-run] would restore %d keys to %s\n", len(snap.Secrets), snap.Path)
		return nil
	}

	payload := make(map[string]interface{}, len(snap.Secrets))
	for k, v := range snap.Secrets {
		payload[k] = v
	}

	_, err := client.Logical().WriteWithContext(ctx, snap.Path, payload)
	if err != nil {
		return fmt.Errorf("rollback: restore %s: %w", snap.Path, err)
	}
	return nil
}
