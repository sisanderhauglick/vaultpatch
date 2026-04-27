// Package truncate provides functionality to shorten secret values
// to a maximum byte length across one or more Vault paths.
package truncate

import (
	"errors"
	"fmt"
	"time"

	"github.com/youorg/vaultpatch/internal/vault"
)

// Options controls how truncation is applied.
type Options struct {
	Paths     []string
	Keys      []string // if empty, all keys are truncated
	MaxLen    int
	Suffix    string // appended when a value is truncated, e.g. "..."
	DryRun    bool
}

// Result holds the outcome for a single path.
type Result struct {
	Path        string
	Truncated   map[string]string // key -> new value
	TruncatedAt time.Time
	DryRun      bool
}

// Truncate reads each path, shortens any value exceeding MaxLen, and
// writes the result back unless DryRun is true.
func Truncate(c *vault.Client, opts Options) ([]Result, error) {
	if c == nil {
		return nil, errors.New("truncate: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("truncate: at least one path is required")
	}
	if opts.MaxLen <= 0 {
		return nil, errors.New("truncate: MaxLen must be greater than zero")
	}

	var results []Result
	for _, p := range opts.Paths {
		secrets, err := vault.ReadSecrets(c, p)
		if err != nil {
			return nil, fmt.Errorf("truncate: read %q: %w", p, err)
		}

		changed := make(map[string]string)
		for k, v := range secrets {
			if !shouldProcess(k, opts.Keys) {
				continue
			}
			s, ok := v.(string)
			if !ok {
				continue
			}
			if len(s) > opts.MaxLen {
				truncated := s[:opts.MaxLen] + opts.Suffix
				changed[k] = truncated
				secrets[k] = truncated
			}
		}

		if !opts.DryRun && len(changed) > 0 {
			if err := vault.WriteSecrets(c, p, secrets); err != nil {
				return nil, fmt.Errorf("truncate: write %q: %w", p, err)
			}
		}

		results = append(results, Result{
			Path:        p,
			Truncated:   changed,
			TruncatedAt: time.Now().UTC(),
			DryRun:      opts.DryRun,
		})
	}
	return results, nil
}

func shouldProcess(key string, keys []string) bool {
	if len(keys) == 0 {
		return true
	}
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}
