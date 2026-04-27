// Package split provides functionality to split a Vault secret path into
// multiple destination paths, distributing keys across targets.
package split

import (
	"errors"
	"fmt"
	"time"

	"github.com/your-org/vaultpatch/internal/vault"
)

// Result holds the outcome of a Split operation for a single destination.
type Result struct {
	Source      string
	Destination string
	Keys        []string
	DryRun      bool
	SplitAt     time.Time
}

// Options controls the behaviour of Split.
type Options struct {
	// Assignments maps destination path -> list of keys to write there.
	Assignments map[string][]string
	DryRun      bool
}

// Split reads secrets from source and writes subsets to each destination
// defined in opts.Assignments. Keys not assigned to any destination are
// silently ignored.
func Split(client *vault.Client, source string, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("split: client must not be nil")
	}
	if source == "" {
		return nil, errors.New("split: source path must not be empty")
	}
	if len(opts.Assignments) == 0 {
		return nil, errors.New("split: assignments must not be empty")
	}

	secrets, err := client.ReadSecrets(source)
	if err != nil {
		return nil, fmt.Errorf("split: read %q: %w", source, err)
	}

	var results []Result
	for dest, keys := range opts.Assignments {
		if dest == "" {
			return nil, errors.New("split: destination path must not be empty")
		}
		subset := make(map[string]interface{}, len(keys))
		for _, k := range keys {
			if v, ok := secrets[k]; ok {
				subset[k] = v
			}
		}
		r := Result{
			Source:      source,
			Destination: dest,
			Keys:        keys,
			DryRun:      opts.DryRun,
			SplitAt:     time.Now().UTC(),
		}
		if !opts.DryRun {
			if err := client.WriteSecrets(dest, subset); err != nil {
				return nil, fmt.Errorf("split: write %q: %w", dest, err)
			}
		}
		results = append(results, r)
	}
	return results, nil
}
