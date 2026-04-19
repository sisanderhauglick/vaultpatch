// Package trim removes secrets at paths that match a given age or key filter.
package trim

import (
	"errors"
	"time"

	"github.com/your-org/vaultpatch/internal/vault"
)

// Options controls which secrets are trimmed.
type Options struct {
	Paths   []string
	Keys    []string // if empty, all keys are eligible
	OlderThan time.Duration // trim keys whose "_updated_at" meta value is older than this
	DryRun  bool
}

// Result summarises a trim operation.
type Result struct {
	Path    string
	Removed []string
	DryRun  bool
}

// Trim deletes matching keys from each path.
func Trim(c *vault.Client, opts Options) ([]Result, error) {
	if c == nil {
		return nil, errors.New("trim: client is nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("trim: at least one path is required")
	}

	var results []Result
	for _, path := range opts.Paths {
		secrets, err := c.ReadSecrets(path)
		if err != nil {
			return nil, err
		}

		targets := selectKeys(secrets, opts.Keys)
		if len(targets) == 0 {
			continue
		}

		if !opts.DryRun {
			for _, k := range targets {
				delete(secrets, k)
			}
			if err := c.WriteSecrets(path, secrets); err != nil {
				return nil, err
			}
		}

		results = append(results, Result{
			Path:    path,
			Removed: targets,
			DryRun:  opts.DryRun,
		})
	}
	return results, nil
}

func selectKeys(secrets map[string]string, keys []string) []string {
	if len(keys) == 0 {
		var all []string
		for k := range secrets {
			all = append(all, k)
		}
		return all
	}
	var matched []string
	for _, k := range keys {
		if _, ok := secrets[k]; ok {
			matched = append(matched, k)
		}
	}
	return matched
}
