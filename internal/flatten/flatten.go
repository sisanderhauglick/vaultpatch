// Package flatten provides functionality to flatten nested Vault secret
// paths into a single destination path, merging all keys.
package flatten

import (
	"errors"
	"fmt"
	"time"

	"github.com/yourusername/vaultpatch/internal/vault"
)

// Result holds the outcome of a Flatten operation for a single destination.
type Result struct {
	Destination string
	Sources     []string
	KeysMerged  int
	DryRun      bool
	FlattenedAt time.Time
}

// Options controls how Flatten behaves.
type Options struct {
	Sources     []string
	Destination string
	Keys        []string // if empty, all keys are included
	Overwrite   bool     // if false, existing keys in dest are preserved
	DryRun      bool
}

// Flatten reads secrets from multiple source paths and writes all keys into a
// single destination path. Conflicts are resolved by the Overwrite option.
func Flatten(client *vault.Client, opts Options) (Result, error) {
	if client == nil {
		return Result{}, errors.New("flatten: client must not be nil")
	}
	if len(opts.Sources) == 0 {
		return Result{}, errors.New("flatten: at least one source path is required")
	}
	if opts.Destination == "" {
		return Result{}, errors.New("flatten: destination path must not be empty")
	}

	merged := make(map[string]interface{})

	// If not overwriting, seed merged with existing destination secrets.
	if !opts.Overwrite {
		existing, err := vault.ReadSecrets(client, opts.Destination)
		if err == nil {
			for k, v := range existing {
				merged[k] = v
			}
		}
	}

	for _, src := range opts.Sources {
		secrets, err := vault.ReadSecrets(client, src)
		if err != nil {
			return Result{}, fmt.Errorf("flatten: read %q: %w", src, err)
		}
		for k, v := range secrets {
			if len(opts.Keys) > 0 && !containsKey(opts.Keys, k) {
				continue
			}
			if _, exists := merged[k]; exists && !opts.Overwrite {
				continue
			}
			merged[k] = v
		}
	}

	result := Result{
		Destination: opts.Destination,
		Sources:     opts.Sources,
		KeysMerged:  len(merged),
		DryRun:      opts.DryRun,
		FlattenedAt: time.Now().UTC(),
	}

	if opts.DryRun {
		return result, nil
	}

	if err := vault.WriteSecrets(client, opts.Destination, merged); err != nil {
		return Result{}, fmt.Errorf("flatten: write %q: %w", opts.Destination, err)
	}

	return result, nil
}

func containsKey(keys []string, target string) bool {
	for _, k := range keys {
		if k == target {
			return true
		}
	}
	return false
}
