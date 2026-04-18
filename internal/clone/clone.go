// Package clone provides functionality to deep-copy a Vault secret path
// tree from one mount/prefix to another.
package clone

import (
	"errors"
	"fmt"

	"github.com/your-org/vaultpatch/internal/vault"
)

// Result holds the outcome of a Clone operation.
type Result struct {
	Source      string
	Destination string
	KeysCopied  int
	DryRun      bool
}

// Options configures a Clone call.
type Options struct {
	Source      string
	Destination string
	Keys        []string // if empty, all keys are cloned
	DryRun      bool
}

// Clone reads secrets from Source and writes them to Destination.
func Clone(client *vault.Client, opts Options) (Result, error) {
	if client == nil {
		return Result{}, errors.New("clone: client must not be nil")
	}
	if opts.Source == "" {
		return Result{}, errors.New("clone: source path must not be empty")
	}
	if opts.Destination == "" {
		return Result{}, errors.New("clone: destination path must not be empty")
	}

	secrets, err := vault.ReadSecrets(client, opts.Source)
	if err != nil {
		return Result{}, fmt.Errorf("clone: read source: %w", err)
	}

	filtered := filterKeys(secrets, opts.Keys)

	if !opts.DryRun {
		if err := vault.WriteSecrets(client, opts.Destination, filtered); err != nil {
			return Result{}, fmt.Errorf("clone: write destination: %w", err)
		}
	}

	return Result{
		Source:      opts.Source,
		Destination: opts.Destination,
		KeysCopied:  len(filtered),
		DryRun:      opts.DryRun,
	}, nil
}

// filterKeys returns a subset of m limited to the provided keys.
// If keys is empty, the original map is returned unchanged.
func filterKeys(m map[string]string, keys []string) map[string]string {
	if len(keys) == 0 {
		return m
	}
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		if v, ok := m[k]; ok {
			out[k] = v
		}
	}
	return out
}
