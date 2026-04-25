// Package migrate moves secrets from one Vault path structure to another,
// optionally remapping key names via a translation map.
package migrate

import (
	"errors"
	"fmt"
	"time"

	"github.com/yourusername/vaultpatch/internal/vault"
)

// Options controls how migration is performed.
type Options struct {
	Sources     []string
	Destination string
	KeyMap      map[string]string // old key -> new key; empty means identity
	DryRun      bool
	Overwrite   bool
}

// Result describes the outcome of migrating a single source path.
type Result struct {
	Source      string
	Destination string
	KeysMapped  int
	DryRun      bool
	MigratedAt  time.Time
}

// Migrate reads each source path, optionally remaps keys, and writes the
// combined result to the destination path.
func Migrate(c *vault.Client, opts Options) ([]Result, error) {
	if c == nil {
		return nil, errors.New("migrate: client must not be nil")
	}
	if len(opts.Sources) == 0 {
		return nil, errors.New("migrate: at least one source path is required")
	}
	if opts.Destination == "" {
		return nil, errors.New("migrate: destination path must not be empty")
	}

	var results []Result

	for _, src := range opts.Sources {
		secrets, err := vault.ReadSecrets(c, src)
		if err != nil {
			return nil, fmt.Errorf("migrate: read %q: %w", src, err)
		}

		mapped := remapKeys(secrets, opts.KeyMap)

		if !opts.DryRun {
			if err := vault.WriteSecrets(c, opts.Destination, mapped); err != nil {
				return nil, fmt.Errorf("migrate: write to %q: %w", opts.Destination, err)
			}
		}

		results = append(results, Result{
			Source:      src,
			Destination: opts.Destination,
			KeysMapped:  len(mapped),
			DryRun:      opts.DryRun,
			MigratedAt:  time.Now().UTC(),
		})
	}

	return results, nil
}

// remapKeys returns a new map with keys renamed according to keyMap.
// Keys absent from keyMap are preserved as-is.
func remapKeys(src map[string]string, keyMap map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		if newKey, ok := keyMap[k]; ok {
			out[newKey] = v
		} else {
			out[k] = v
		}
	}
	return out
}
