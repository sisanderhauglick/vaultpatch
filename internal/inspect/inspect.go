// Package inspect provides metadata inspection for Vault secret paths.
package inspect

import (
	"errors"
	"time"

	"github.com/youorg/vaultpatch/internal/vault"
)

// Result holds metadata about a single secret path.
type Result struct {
	Path      string
	KeyCount  int
	Keys      []string
	FetchedAt time.Time
}

// Options controls Inspect behaviour.
type Options struct {
	Paths []string
}

// Inspect reads metadata for each path and returns a slice of Results.
func Inspect(client *vault.Client, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("inspect: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("inspect: at least one path is required")
	}

	var results []Result
	for _, path := range opts.Paths {
		secrets, err := vault.ReadSecrets(client, path)
		if err != nil {
			return nil, err
		}

		keys := make([]string, 0, len(secrets))
		for k := range secrets {
			keys = append(keys, k)
		}

		results = append(results, Result{
			Path:      path,
			KeyCount:  len(keys),
			Keys:      keys,
			FetchedAt: time.Now().UTC(),
		})
	}
	return results, nil
}
