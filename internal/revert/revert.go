// Package revert provides functionality to revert specific keys at a Vault
// path back to a previous value captured in a before-map.
package revert

import (
	"errors"
	"fmt"
	"time"

	"github.com/user/vaultpatch/internal/vault"
)

// Result describes the outcome of a revert operation for a single path.
type Result struct {
	Path      string
	Reverted  []string
	Skipped   []string
	DryRun    bool
	RevertedAt time.Time
}

// Options controls Revert behaviour.
type Options struct {
	Paths  []string
	Before map[string]string // key -> previous value
	Keys   []string          // if empty, revert all keys present in Before
	DryRun bool
}

// Revert writes previous values back to Vault for the specified paths and keys.
func Revert(client *vault.Client, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("revert: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("revert: at least one path is required")
	}
	if len(opts.Before) == 0 {
		return nil, errors.New("revert: before map must not be empty")
	}

	targetKeys := opts.Keys
	if len(targetKeys) == 0 {
		for k := range opts.Before {
			targetKeys = append(targetKeys, k)
		}
	}

	var results []Result
	for _, path := range opts.Paths {
		res := Result{Path: path, DryRun: opts.DryRun}

		current, err := vault.ReadSecrets(client, path)
		if err != nil {
			return nil, fmt.Errorf("revert: read %s: %w", path, err)
		}

		patch := make(map[string]interface{})
		for _, k := range targetKeys {
			prev, ok := opts.Before[k]
			if !ok {
				res.Skipped = append(res.Skipped, k)
				continue
			}
			if cur, exists := current[k]; exists && cur == prev {
				res.Skipped = append(res.Skipped, k)
				continue
			}
			patch[k] = prev
			res.Reverted = append(res.Reverted, k)
		}

		if !opts.DryRun && len(patch) > 0 {
			for k, v := range current {
				if _, overriding := patch[k]; !overriding {
					patch[k] = v
				}
			}
			if err := vault.WriteSecrets(client, path, patch); err != nil {
				return nil, fmt.Errorf("revert: write %s: %w", path, err)
			}
			res.RevertedAt = time.Now().UTC()
		}

		results = append(results, res)
	}
	return results, nil
}
