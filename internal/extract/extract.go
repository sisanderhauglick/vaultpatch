// Package extract provides functionality to extract a subset of keys
// from one or more Vault secret paths into a new destination path.
package extract

import (
	"errors"
	"fmt"
	"time"

	"github.com/youorg/vaultpatch/internal/vault"
)

// Result holds the outcome of an Extract operation for a single source path.
type Result struct {
	Source      string
	Destination string
	Keys        []string
	ExtractedAt time.Time
	DryRun      bool
}

// Options configures the Extract operation.
type Options struct {
	Sources     []string
	Destination string
	Keys        []string
	DryRun      bool
}

// Extract reads the specified keys from each source path and writes them
// to the destination path. If Keys is empty, all keys are extracted.
func Extract(c *vault.Client, opts Options) ([]Result, error) {
	if c == nil {
		return nil, errors.New("extract: client must not be nil")
	}
	if len(opts.Sources) == 0 {
		return nil, errors.New("extract: at least one source path is required")
	}
	if opts.Destination == "" {
		return nil, errors.New("extract: destination path must not be empty")
	}

	merged := make(map[string]string)
	var results []Result

	for _, src := range opts.Sources {
		data, err := vault.ReadSecrets(c, src)
		if err != nil {
			return nil, fmt.Errorf("extract: read %q: %w", src, err)
		}

		picked := pickKeys(data, opts.Keys)
		for k, v := range picked {
			merged[k] = v
		}

		results = append(results, Result{
			Source:      src,
			Destination: opts.Destination,
			Keys:        keyList(picked),
			ExtractedAt: time.Now().UTC(),
			DryRun:      opts.DryRun,
		})
	}

	if !opts.DryRun && len(merged) > 0 {
		if err := vault.WriteSecrets(c, opts.Destination, merged); err != nil {
			return nil, fmt.Errorf("extract: write %q: %w", opts.Destination, err)
		}
	}

	return results, nil
}

// pickKeys returns the subset of data matching the requested keys.
// If keys is empty, all entries are returned.
func pickKeys(data map[string]string, keys []string) map[string]string {
	if len(keys) == 0 {
		out := make(map[string]string, len(data))
		for k, v := range data {
			out[k] = v
		}
		return out
	}
	out := make(map[string]string)
	for _, k := range keys {
		if v, ok := data[k]; ok {
			out[k] = v
		}
	}
	return out
}

func keyList(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
