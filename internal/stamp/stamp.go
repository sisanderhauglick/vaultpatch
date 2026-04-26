// Package stamp applies metadata annotations (stamps) to Vault secret paths.
// Stamps are written as special keys (e.g. __stamped_by, __stamped_at) alongside
// existing secret data, enabling traceability across environments.
package stamp

import (
	"errors"
	"fmt"
	"time"

	"github.com/youorg/vaultpatch/internal/vault"
)

// Result holds the outcome of a Stamp operation for a single path.
type Result struct {
	Path       string
	Stamped    bool
	DryRun     bool
	StampedAt  time.Time
	Annotations map[string]string
}

// Options configures a Stamp call.
type Options struct {
	Paths       []string
	Annotations map[string]string
	DryRun      bool
}

// Stamp writes annotation keys into each secret at the given paths.
// If DryRun is true, no writes are performed but results are still returned.
func Stamp(client *vault.Client, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("stamp: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("stamp: at least one path is required")
	}
	if len(opts.Annotations) == 0 {
		return nil, errors.New("stamp: at least one annotation is required")
	}

	now := time.Now().UTC()
	var results []Result

	for _, path := range opts.Paths {
		existing, err := vault.ReadSecrets(client, path)
		if err != nil {
			return nil, fmt.Errorf("stamp: read %q: %w", path, err)
		}

		merged := make(map[string]interface{}, len(existing)+len(opts.Annotations))
		for k, v := range existing {
			merged[k] = v
		}
		for k, v := range opts.Annotations {
			merged["__"+k] = v
		}
		merged["__stamped_at"] = now.Format(time.RFC3339)

		if !opts.DryRun {
			if err := vault.WriteSecrets(client, path, merged); err != nil {
				return nil, fmt.Errorf("stamp: write %q: %w", path, err)
			}
		}

		results = append(results, Result{
			Path:        path,
			Stamped:     !opts.DryRun,
			DryRun:      opts.DryRun,
			StampedAt:   now,
			Annotations: opts.Annotations,
		})
	}

	return results, nil
}
