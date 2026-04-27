// Package prefill populates missing keys in Vault secret paths using a
// set of default values, leaving existing keys untouched.
package prefill

import (
	"errors"
	"fmt"
	"time"

	"github.com/youorg/vaultpatch/internal/vault"
)

// Options controls Prefill behaviour.
type Options struct {
	// Paths are the Vault KV paths to prefill.
	Paths []string
	// Defaults is the map of key→value pairs to inject when absent.
	Defaults map[string]string
	// DryRun skips writing changes back to Vault.
	DryRun bool
}

// Result describes what happened for a single path.
type Result struct {
	Path        string
	Filled      map[string]string
	DryRun      bool
	PrefillledAt time.Time
}

// Prefill reads each path, identifies missing default keys, and writes
// them back unless DryRun is set.
func Prefill(c *vault.Client, opts Options) ([]Result, error) {
	if c == nil {
		return nil, errors.New("prefill: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("prefill: at least one path is required")
	}
	if len(opts.Defaults) == 0 {
		return nil, errors.New("prefill: defaults map must not be empty")
	}

	var results []Result
	for _, p := range opts.Paths {
		existing, err := vault.ReadSecrets(c, p)
		if err != nil {
			return nil, fmt.Errorf("prefill: read %q: %w", p, err)
		}

		filled := make(map[string]string)
		for k, v := range opts.Defaults {
			if _, ok := existing[k]; !ok {
				filled[k] = v
			}
		}

		if !opts.DryRun && len(filled) > 0 {
			merged := make(map[string]interface{}, len(existing)+len(filled))
			for k, v := range existing {
				merged[k] = v
			}
			for k, v := range filled {
				merged[k] = v
			}
			if err := vault.WriteSecrets(c, p, merged); err != nil {
				return nil, fmt.Errorf("prefill: write %q: %w", p, err)
			}
		}

		results = append(results, Result{
			Path:        p,
			Filled:      filled,
			DryRun:      opts.DryRun,
			PrefillledAt: time.Now().UTC(),
		})
	}
	return results, nil
}
