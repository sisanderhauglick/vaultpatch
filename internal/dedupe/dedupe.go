// Package dedupe removes duplicate keys across multiple Vault secret paths,
// keeping the value from the highest-priority (first) source.
package dedupe

import (
	"errors"
	"time"

	"github.com/yourusername/vaultpatch/internal/vault"
)

// Result holds the outcome of a deduplication operation for a single path.
type Result struct {
	Path        string
	RemovedKeys []string
	DryRun      bool
	DedupedAt   time.Time
}

// Options configures a Dedupe run.
type Options struct {
	// Paths is the ordered list of secret paths; earlier paths take priority.
	Paths  []string
	DryRun bool
}

// Dedupe reads secrets from all paths and removes keys that already appear in
// a higher-priority path. When DryRun is true no writes are performed.
func Dedupe(client *vault.Client, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("dedupe: client must not be nil")
	}
	if len(opts.Paths) < 2 {
		return nil, errors.New("dedupe: at least two paths are required")
	}

	// Build a set of keys already seen in higher-priority paths.
	seen := make(map[string]bool)
	var results []Result

	for _, path := range opts.Paths {
		secrets, err := vault.ReadSecrets(client, path)
		if err != nil {
			return nil, err
		}

		var removed []string
		clean := make(map[string]interface{})

		for k, v := range secrets {
			if seen[k] {
				removed = append(removed, k)
			} else {
				clean[k] = v
				seen[k] = true
			}
		}

		if !opts.DryRun && len(removed) > 0 {
			if err := vault.WriteSecrets(client, path, clean); err != nil {
				return nil, err
			}
		}

		results = append(results, Result{
			Path:        path,
			RemovedKeys: removed,
			DryRun:      opts.DryRun,
			DedupedAt:   time.Now().UTC(),
		})
	}

	return results, nil
}
