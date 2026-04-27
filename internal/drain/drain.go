// Package drain removes all secrets from one or more Vault paths,
// optionally preserving a set of keys and supporting dry-run mode.
package drain

import (
	"errors"
	"fmt"
	"time"

	"github.com/your-org/vaultpatch/internal/vault"
)

// Result describes the outcome of a Drain operation on a single path.
type Result struct {
	Path      string
	Drained   []string
	Preserved []string
	DrainedAt time.Time
	DryRun    bool
}

// Options controls Drain behaviour.
type Options struct {
	// Paths are the Vault KV paths to drain.
	Paths []string
	// Preserve lists keys that must not be deleted.
	Preserve []string
	// DryRun reports what would be removed without writing.
	DryRun bool
}

// Drain removes every key from each path except those listed in Options.Preserve.
// When DryRun is true no writes are performed.
func Drain(client *vault.Client, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("drain: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("drain: at least one path is required")
	}

	preserveSet := make(map[string]struct{}, len(opts.Preserve))
	for _, k := range opts.Preserve {
		preserveSet[k] = struct{}{}
	}

	var results []Result
	for _, path := range opts.Paths {
		secrets, err := client.Read(path)
		if err != nil {
			return nil, fmt.Errorf("drain: read %q: %w", path, err)
		}

		var drained, preserved []string
		filtered := make(map[string]interface{})
		for k, v := range secrets {
			if _, keep := preserveSet[k]; keep {
				filtered[k] = v
				preserved = append(preserved, k)
			} else {
				drained = append(drained, k)
			}
		}

		if !opts.DryRun && len(drained) > 0 {
			if err := client.Write(path, filtered); err != nil {
				return nil, fmt.Errorf("drain: write %q: %w", path, err)
			}
		}

		results = append(results, Result{
			Path:      path,
			Drained:   drained,
			Preserved: preserved,
			DrainedAt: time.Now().UTC(),
			DryRun:    opts.DryRun,
		})
	}
	return results, nil
}
