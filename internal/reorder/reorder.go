package reorder

import (
	"errors"
	"fmt"
	"time"

	"github.com/vaultpatch/vaultpatch/internal/vault"
)

// Result holds the outcome of a Reorder operation for a single path.
type Result struct {
	Path      string
	OldOrder  []string
	NewOrder  []string
	ReorderedAt time.Time
	DryRun    bool
}

// Options controls Reorder behaviour.
type Options struct {
	// Keys defines the desired key order. Keys not listed are appended in
	// their original relative order after the explicitly ordered keys.
	Keys   []string
	DryRun bool
}

// Reorder reads secrets at each path and writes them back with keys in the
// order defined by opts.Keys. Keys absent from opts.Keys are preserved at
// the end in their original order.
func Reorder(c *vault.Client, paths []string, opts Options) ([]Result, error) {
	if c == nil {
		return nil, errors.New("reorder: client must not be nil")
	}
	if len(paths) == 0 {
		return nil, errors.New("reorder: at least one path is required")
	}
	if len(opts.Keys) == 0 {
		return nil, errors.New("reorder: at least one key must be specified")
	}

	var results []Result
	for _, p := range paths {
		secrets, err := vault.ReadSecrets(c, p)
		if err != nil {
			return nil, fmt.Errorf("reorder: read %q: %w", p, err)
		}

		oldOrder := sortedKeys(secrets)
		newOrder := buildOrder(opts.Keys, oldOrder)

		r := Result{
			Path:        p,
			OldOrder:    oldOrder,
			NewOrder:    newOrder,
			ReorderedAt: time.Now().UTC(),
			DryRun:      opts.DryRun,
		}

		if !opts.DryRun {
			ordered := make(map[string]interface{}, len(secrets))
			for k, v := range secrets {
				ordered[k] = v
			}
			if err := vault.WriteSecrets(c, p, ordered); err != nil {
				return nil, fmt.Errorf("reorder: write %q: %w", p, err)
			}
		}

		results = append(results, r)
	}
	return results, nil
}

// buildOrder returns the desired key ordering: explicit keys first (in the
// order given), then remaining keys in their original order.
func buildOrder(explicit, original []string) []string {
	seen := make(map[string]bool, len(explicit))
	for _, k := range explicit {
		seen[k] = true
	}

	ordered := make([]string, 0, len(original))
	for _, k := range explicit {
		for _, ok := range original {
			if ok == k {
				ordered = append(ordered, k)
				break
			}
		}
	}
	for _, k := range original {
		if !seen[k] {
			ordered = append(ordered, k)
		}
	}
	return ordered
}

// sortedKeys returns the map keys in their natural iteration order (used to
// capture a stable snapshot of the existing order).
func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
