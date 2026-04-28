// Package enrich adds computed metadata fields to Vault secrets at given paths.
package enrich

import (
	"errors"
	"fmt"
	"time"
)

// VaultClient is the subset of the Vault client used by Enrich.
type VaultClient interface {
	Read(path string) (map[string]interface{}, error)
	Write(path string, data map[string]interface{}) error
}

// Options controls Enrich behaviour.
type Options struct {
	Paths       []string
	Annotations map[string]string // static key→value pairs to inject
	DryRun      bool
}

// Result holds the outcome for a single path.
type Result struct {
	Path       string
	EnrichedAt time.Time
	Added      []string
	DryRun     bool
}

// Enrich reads each path, merges the provided annotations, and writes back.
// On dry-run the write is skipped but results are still populated.
func Enrich(client VaultClient, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("enrich: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("enrich: at least one path is required")
	}
	if len(opts.Annotations) == 0 {
		return nil, errors.New("enrich: at least one annotation is required")
	}

	var results []Result
	for _, path := range opts.Paths {
		data, err := client.Read(path)
		if err != nil {
			return nil, fmt.Errorf("enrich: read %q: %w", path, err)
		}
		if data == nil {
			data = make(map[string]interface{})
		}

		var added []string
		for k, v := range opts.Annotations {
			data[k] = v
			added = append(added, k)
		}

		if !opts.DryRun {
			if err := client.Write(path, data); err != nil {
				return nil, fmt.Errorf("enrich: write %q: %w", path, err)
			}
		}

		results = append(results, Result{
			Path:       path,
			EnrichedAt: time.Now().UTC(),
			Added:      added,
			DryRun:     opts.DryRun,
		})
	}
	return results, nil
}
