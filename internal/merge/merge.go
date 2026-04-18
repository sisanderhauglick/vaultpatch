// Package merge provides functionality to merge secrets from multiple
// Vault paths into a single destination path.
package merge

import (
	"errors"
	"fmt"

	"github.com/your-org/vaultpatch/internal/vault"
)

// Result holds the outcome of a merge operation.
type Result struct {
	Destination string
	Merged      int
	DryRun      bool
}

// Options configures a Merge call.
type Options struct {
	Sources     []string
	Destination string
	// Keys restricts which keys are merged; empty means all keys.
	Keys   []string
	DryRun bool
}

// Merge reads secrets from all source paths and writes the combined map to
// the destination path. Later sources take precedence over earlier ones.
func Merge(client *vault.Client, opts Options) (Result, error) {
	if client == nil {
		return Result{}, errors.New("merge: client must not be nil")
	}
	if len(opts.Sources) == 0 {
		return Result{}, errors.New("merge: at least one source path is required")
	}
	if opts.Destination == "" {
		return Result{}, errors.New("merge: destination path must not be empty")
	}

	combined := make(map[string]string)

	for _, src := range opts.Sources {
		data, err := vault.ReadSecrets(client, src)
		if err != nil {
			return Result{}, fmt.Errorf("merge: reading %q: %w", src, err)
		}
		for k, v := range data {
			if len(opts.Keys) == 0 || containsKey(opts.Keys, k) {
				combined[k] = v
			}
		}
	}

	if opts.DryRun {
		return Result{Destination: opts.Destination, Merged: len(combined), DryRun: true}, nil
	}

	if err := vault.WriteSecrets(client, opts.Destination, combined); err != nil {
		return Result{}, fmt.Errorf("merge: writing to %q: %w", opts.Destination, err)
	}

	return Result{Destination: opts.Destination, Merged: len(combined), DryRun: false}, nil
}

func containsKey(keys []string, k string) bool {
	for _, key := range keys {
		if key == k {
			return true
		}
	}
	return false
}
