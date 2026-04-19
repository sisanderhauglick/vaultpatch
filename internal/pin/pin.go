// Package pin provides functionality to pin Vault secrets at a specific
// version, preventing unintended overwrites during promotions or syncs.
package pin

import (
	"errors"
	"fmt"
	"time"

	"github.com/your-org/vaultpatch/internal/vault"
)

// Result holds the outcome of a pin or unpin operation for a single path.
type Result struct {
	Path     string
	Version  int
	Pinned   bool
	DryRun   bool
	PinnedAt time.Time
}

// Options configures a Pin or Unpin call.
type Options struct {
	Paths   []string
	Version int
	DryRun  bool
}

const pinMetaKey = "_vaultpatch_pinned_version"

// Pin marks each path's metadata with the target version so downstream
// operations can detect and respect the pin.
func Pin(client *vault.Client, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("pin: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("pin: at least one path is required")
	}
	if opts.Version < 1 {
		return nil, errors.New("pin: version must be >= 1")
	}

	var results []Result
	for _, path := range opts.Paths {
		r := Result{
			Path:     path,
			Version:  opts.Version,
			Pinned:   true,
			DryRun:   opts.DryRun,
			PinnedAt: time.Now().UTC(),
		}
		if !opts.DryRun {
			existing, err := client.Read(path)
			if err != nil {
				return nil, fmt.Errorf("pin: read %s: %w", path, err)
			}
			if existing == nil {
				existing = make(map[string]string)
			}
			existing[pinMetaKey] = fmt.Sprintf("%d", opts.Version)
			if err := client.Write(path, existing); err != nil {
				return nil, fmt.Errorf("pin: write %s: %w", path, err)
			}
		}
		results = append(results, r)
	}
	return results, nil
}

// Unpin removes the pin metadata from each path.
func Unpin(client *vault.Client, paths []string, dryRun bool) ([]Result, error) {
	if client == nil {
		return nil, errors.New("unpin: client must not be nil")
	}
	if len(paths) == 0 {
		return nil, errors.New("unpin: at least one path is required")
	}

	var results []Result
	for _, path := range paths {
		r := Result{Path: path, Pinned: false, DryRun: dryRun}
		if !dryRun {
			existing, err := client.Read(path)
			if err != nil {
				return nil, fmt.Errorf("unpin: read %s: %w", path, err)
			}
			if existing != nil {
				delete(existing, pinMetaKey)
				if err := client.Write(path, existing); err != nil {
					return nil, fmt.Errorf("unpin: write %s: %w", path, err)
				}
			}
		}
		results = append(results, r)
	}
	return results, nil
}
