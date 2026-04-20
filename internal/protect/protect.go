// Package protect provides secret path protection (write-lock) functionality.
// Protected paths reject write operations until explicitly unprotected.
package protect

import (
	"errors"
	"fmt"
	"time"

	"github.com/your-org/vaultpatch/internal/vault"
)

const protectMetaKey = "__vaultpatch_protected"

// Result holds the outcome of a Protect or Unprotect operation on a single path.
type Result struct {
	Path        string
	Protected   bool
	DryRun      bool
	ProtectedAt time.Time
	Owner       string
}

// Options configures a Protect or Unprotect call.
type Options struct {
	Paths  []string
	Owner  string
	DryRun bool
}

// Protect marks each path as protected by writing metadata.
func Protect(client *vault.Client, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("protect: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("protect: at least one path is required")
	}
	if opts.Owner == "" {
		return nil, errors.New("protect: owner must not be empty")
	}

	now := time.Now().UTC()
	var results []Result

	for _, path := range opts.Paths {
		r := Result{
			Path:        path,
			Protected:   true,
			DryRun:      opts.DryRun,
			ProtectedAt: now,
			Owner:       opts.Owner,
		}
		if !opts.DryRun {
			existing, err := vault.ReadSecrets(client, path)
			if err != nil {
				return nil, fmt.Errorf("protect: read %q: %w", path, err)
			}
			existing[protectMetaKey] = fmt.Sprintf("%s|%s", opts.Owner, now.Format(time.RFC3339))
			if err := vault.WriteSecrets(client, path, existing); err != nil {
				return nil, fmt.Errorf("protect: write %q: %w", path, err)
			}
		}
		results = append(results, r)
	}
	return results, nil
}

// Unprotect removes the protection marker from each path.
func Unprotect(client *vault.Client, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("protect: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("protect: at least one path is required")
	}

	var results []Result
	for _, path := range opts.Paths {
		r := Result{Path: path, Protected: false, DryRun: opts.DryRun}
		if !opts.DryRun {
			existing, err := vault.ReadSecrets(client, path)
			if err != nil {
				return nil, fmt.Errorf("protect: read %q: %w", path, err)
			}
			delete(existing, protectMetaKey)
			if err := vault.WriteSecrets(client, path, existing); err != nil {
				return nil, fmt.Errorf("protect: write %q: %w", path, err)
			}
		}
		results = append(results, r)
	}
	return results, nil
}

// IsProtected reports whether a path currently carries the protection marker.
func IsProtected(client *vault.Client, path string) (bool, error) {
	if client == nil {
		return false, errors.New("protect: client must not be nil")
	}
	secrets, err := vault.ReadSecrets(client, path)
	if err != nil {
		return false, fmt.Errorf("protect: read %q: %w", path, err)
	}
	_, ok := secrets[protectMetaKey]
	return ok, nil
}
