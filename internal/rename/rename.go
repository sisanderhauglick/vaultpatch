// Package rename provides functionality to rename secret keys across Vault paths.
package rename

import (
	"errors"
	"fmt"

	"vaultpatch/internal/vault"
)

// Result holds the outcome of a rename operation.
type Result struct {
	Path    string
	OldKey  string
	NewKey  string
	DryRun  bool
	Skipped bool
}

// Options configures the rename operation.
type Options struct {
	Paths  []string
	OldKey string
	NewKey string
	DryRun bool
}

// Rename reads each path, copies the value from OldKey to NewKey, and removes OldKey.
// In dry-run mode no writes are performed.
func Rename(client *vault.Client, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("rename: vault client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("rename: at least one path is required")
	}
	if opts.OldKey == "" || opts.NewKey == "" {
		return nil, errors.New("rename: old-key and new-key must not be empty")
	}

	var results []Result

	for _, path := range opts.Paths {
		secrets, err := vault.ReadSecrets(client, path)
		if err != nil {
			return nil, fmt.Errorf("rename: read %s: %w", path, err)
		}

		val, ok := secrets[opts.OldKey]
		if !ok {
			results = append(results, Result{Path: path, OldKey: opts.OldKey, NewKey: opts.NewKey, DryRun: opts.DryRun, Skipped: true})
			continue
		}

		if !opts.DryRun {
			updated := make(map[string]string, len(secrets))
			for k, v := range secrets {
				updated[k] = v
			}
			updated[opts.NewKey] = val
			delete(updated, opts.OldKey)

			if err := vault.WriteSecrets(client, path, updated); err != nil {
				return nil, fmt.Errorf("rename: write %s: %w", path, err)
			}
		}

		results = append(results, Result{Path: path, OldKey: opts.OldKey, NewKey: opts.NewKey, DryRun: opts.DryRun, Skipped: false})
	}

	return results, nil
}
