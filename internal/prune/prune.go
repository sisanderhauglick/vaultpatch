// Package prune removes secrets from Vault paths that match an age or version threshold.
package prune

import (
	"errors"
	"fmt"
	"time"
)

// Client is the subset of vault.Client used by Prune.
type Client interface {
	Read(path string) (map[string]interface{}, error)
	Write(path string, data map[string]interface{}) error
	Delete(path string) error
}

// Options controls Prune behaviour.
type Options struct {
	Paths     []string
	OlderThan time.Duration // prune entries whose "_created_at" is older than this
	DryRun    bool
}

// Result describes what happened to a single path.
type Result struct {
	Path    string
	Pruned  bool
	DryRun  bool
	PrunedAt time.Time
	Reason  string
}

// Prune inspects each path and removes secrets that exceed the age threshold.
func Prune(client Client, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("prune: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("prune: at least one path is required")
	}
	if opts.OlderThan <= 0 {
		return nil, errors.New("prune: OlderThan duration must be positive")
	}

	cutoff := time.Now().UTC().Add(-opts.OlderThan)
	var results []Result

	for _, path := range opts.Paths {
		secrets, err := client.Read(path)
		if err != nil {
			return nil, fmt.Errorf("prune: read %q: %w", path, err)
		}

		createdRaw, ok := secrets["_created_at"]
		if !ok {
			results = append(results, Result{
				Path:   path,
				Pruned: false,
				DryRun: opts.DryRun,
				Reason: "no _created_at key",
			})
			continue
		}

		createdStr, _ := createdRaw.(string)
		createdAt, err := time.Parse(time.RFC3339, createdStr)
		if err != nil {
			return nil, fmt.Errorf("prune: parse _created_at for %q: %w", path, err)
		}

		if createdAt.After(cutoff) {
			results = append(results, Result{
				Path:   path,
				Pruned: false,
				DryRun: opts.DryRun,
				Reason: "not old enough",
			})
			continue
		}

		if !opts.DryRun {
			if err := client.Delete(path); err != nil {
				return nil, fmt.Errorf("prune: delete %q: %w", path, err)
			}
		}

		results = append(results, Result{
			Path:     path,
			Pruned:   true,
			DryRun:   opts.DryRun,
			PrunedAt: time.Now().UTC(),
			Reason:   "exceeded age threshold",
		})
	}

	return results, nil
}
