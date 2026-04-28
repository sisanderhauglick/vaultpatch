// Package group provides functionality to group secrets from multiple Vault
// paths under a single composite key namespace.
package group

import (
	"errors"
	"fmt"
	"time"
)

// VaultClient is the interface required by Group.
type VaultClient interface {
	Read(path string) (map[string]interface{}, error)
	Write(path string, data map[string]interface{}) error
}

// Options configures a Group operation.
type Options struct {
	Sources     []string
	Destination string
	Prefix      bool // if true, prefix each key with its source path segment
	DryRun      bool
}

// Result describes the outcome of a Group operation.
type Result struct {
	Destination string
	KeysMerged  int
	GroupedAt   time.Time
	DryRun      bool
}

// Group reads secrets from each source path and writes them combined into
// destination. Key collisions are resolved by last-writer-wins (source order).
func Group(client VaultClient, opts Options) (Result, error) {
	if client == nil {
		return Result{}, errors.New("group: client must not be nil")
	}
	if len(opts.Sources) == 0 {
		return Result{}, errors.New("group: at least one source path is required")
	}
	if opts.Destination == "" {
		return Result{}, errors.New("group: destination path must not be empty")
	}

	merged := make(map[string]interface{})

	for _, src := range opts.Sources {
		data, err := client.Read(src)
		if err != nil {
			return Result{}, fmt.Errorf("group: read %q: %w", src, err)
		}
		for k, v := range data {
			key := k
			if opts.Prefix {
				key = prefixKey(src, k)
			}
			merged[key] = v
		}
	}

	if !opts.DryRun {
		if err := client.Write(opts.Destination, merged); err != nil {
			return Result{}, fmt.Errorf("group: write %q: %w", opts.Destination, err)
		}
	}

	return Result{
		Destination: opts.Destination,
		KeysMerged:  len(merged),
		GroupedAt:   time.Now().UTC(),
		DryRun:      opts.DryRun,
	}, nil
}

// prefixKey returns "<last segment of path>/<key>".
func prefixKey(path, key string) string {
	seg := path
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			seg = path[i+1:]
			break
		}
	}
	return seg + "/" + key
}
