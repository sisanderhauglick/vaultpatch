// Package broadcast copies a secret from a single source path to multiple
// destination paths in Vault, optionally filtering which keys are propagated.
package broadcast

import (
	"errors"
	"fmt"
	"time"

	"github.com/your-org/vaultpatch/internal/vault"
)

// Result describes the outcome of a single broadcast operation.
type Result struct {
	Source      string
	Destination string
	KeysCopied  int
	DryRun      bool
	BroadcastAt time.Time
}

// Options controls broadcast behaviour.
type Options struct {
	Keys   []string // if empty, all keys are broadcast
	DryRun bool
}

// Broadcast reads the secret at source and writes it to every path in dests.
func Broadcast(c *vault.Client, source string, dests []string, opts Options) ([]Result, error) {
	if c == nil {
		return nil, errors.New("broadcast: client is nil")
	}
	if source == "" {
		return nil, errors.New("broadcast: source path is empty")
	}
	if len(dests) == 0 {
		return nil, errors.New("broadcast: destination list is empty")
	}

	secrets, err := vault.ReadSecrets(c, source)
	if err != nil {
		return nil, fmt.Errorf("broadcast: read %s: %w", source, err)
	}

	payload := filterKeys(secrets, opts.Keys)

	var results []Result
	for _, dest := range dests {
		if dest == "" {
			continue
		}
		r := Result{
			Source:      source,
			Destination: dest,
			KeysCopied:  len(payload),
			DryRun:      opts.DryRun,
			BroadcastAt: time.Now().UTC(),
		}
		if !opts.DryRun {
			if werr := vault.WriteSecrets(c, dest, payload); werr != nil {
				return nil, fmt.Errorf("broadcast: write %s: %w", dest, werr)
			}
		}
		results = append(results, r)
	}
	return results, nil
}

// filterKeys returns a subset of data containing only the requested keys.
// If keys is empty the full map is returned unchanged.
func filterKeys(data map[string]string, keys []string) map[string]string {
	if len(keys) == 0 {
		out := make(map[string]string, len(data))
		for k, v := range data {
			out[k] = v
		}
		return out
	}
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		if v, ok := data[k]; ok {
			out[k] = v
		}
	}
	return out
}
