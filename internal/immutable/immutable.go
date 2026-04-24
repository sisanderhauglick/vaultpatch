// Package immutable provides functionality to mark Vault secret paths as
// immutable, preventing any further writes or modifications to those paths.
package immutable

import (
	"errors"
	"fmt"
	"time"
)

const immutableKey = "__immutable"

// VaultClient is the interface required by immutable operations.
type VaultClient interface {
	Read(path string) (map[string]interface{}, error)
	Write(path string, data map[string]interface{}) error
}

// Result holds the outcome of a single path's immutable operation.
type Result struct {
	Path        string
	Immutable   bool
	Released    bool
	DryRun      bool
	LockedAt    time.Time
}

// Options configures the Immute and Release operations.
type Options struct {
	Paths  []string
	DryRun bool
}

// Immute marks the given paths as immutable by writing a sentinel key.
func Immute(client VaultClient, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("immutable: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("immutable: at least one path is required")
	}

	var results []Result
	for _, path := range opts.Paths {
		existing, err := client.Read(path)
		if err != nil {
			return nil, fmt.Errorf("immutable: read %q: %w", path, err)
		}
		if existing == nil {
			existing = make(map[string]interface{})
		}

		now := time.Now().UTC()
		r := Result{Path: path, Immutable: true, DryRun: opts.DryRun, LockedAt: now}

		if !opts.DryRun {
			existing[immutableKey] = now.Format(time.RFC3339)
			if err := client.Write(path, existing); err != nil {
				return nil, fmt.Errorf("immutable: write %q: %w", path, err)
			}
		}
		results = append(results, r)
	}
	return results, nil
}

// Release removes the immutable sentinel key from the given paths.
func Release(client VaultClient, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("immutable: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("immutable: at least one path is required")
	}

	var results []Result
	for _, path := range opts.Paths {
		existing, err := client.Read(path)
		if err != nil {
			return nil, fmt.Errorf("immutable: read %q: %w", path, err)
		}
		if existing == nil {
			existing = make(map[string]interface{})
		}

		r := Result{Path: path, Released: true, DryRun: opts.DryRun}

		if !opts.DryRun {
			delete(existing, immutableKey)
			if err := client.Write(path, existing); err != nil {
				return nil, fmt.Errorf("immutable: write %q: %w", path, err)
			}
		}
		results = append(results, r)
	}
	return results, nil
}

// IsImmutable returns true if the path currently carries the immutable sentinel.
func IsImmutable(client VaultClient, path string) (bool, error) {
	if client == nil {
		return false, errors.New("immutable: client must not be nil")
	}
	data, err := client.Read(path)
	if err != nil {
		return false, fmt.Errorf("immutable: read %q: %w", path, err)
	}
	_, ok := data[immutableKey]
	return ok, nil
}
