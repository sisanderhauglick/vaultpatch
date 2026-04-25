// Package inherit provides functionality to propagate secrets from a parent
// Vault path to one or more child paths, merging keys without overwriting
// existing values unless explicitly forced.
package inherit

import (
	"errors"
	"fmt"
	"time"
)

// VaultClient defines the subset of Vault operations required by Inherit.
type VaultClient interface {
	Read(path string) (map[string]string, error)
	Write(path string, data map[string]string) error
}

// Options controls the behaviour of an Inherit run.
type Options struct {
	// Parent is the source path whose secrets are propagated.
	Parent string
	// Children are the destination paths that inherit from Parent.
	Children []string
	// Keys limits propagation to specific keys; empty means all keys.
	Keys []string
	// Force overwrites keys that already exist in child paths.
	Force bool
	// DryRun skips writes and only reports what would change.
	DryRun bool
}

// Result describes the outcome for a single child path.
type Result struct {
	Path        string
	Inherited   []string
	Skipped     []string
	DryRun      bool
	InheritedAt time.Time
}

// Inherit reads secrets from the parent path and propagates them to each
// child path according to the supplied Options.
func Inherit(client VaultClient, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("inherit: client must not be nil")
	}
	if opts.Parent == "" {
		return nil, errors.New("inherit: parent path must not be empty")
	}
	if len(opts.Children) == 0 {
		return nil, errors.New("inherit: at least one child path is required")
	}

	parentSecrets, err := client.Read(opts.Parent)
	if err != nil {
		return nil, fmt.Errorf("inherit: reading parent %q: %w", opts.Parent, err)
	}

	source := filterKeys(parentSecrets, opts.Keys)

	var results []Result
	for _, child := range opts.Children {
		res, err := propagate(client, child, source, opts)
		if err != nil {
			return nil, fmt.Errorf("inherit: propagating to %q: %w", child, err)
		}
		results = append(results, res)
	}
	return results, nil
}

func propagate(client VaultClient, path string, source map[string]string, opts Options) (Result, error) {
	existing, err := client.Read(path)
	if err != nil {
		return Result{}, err
	}
	if existing == nil {
		existing = map[string]string{}
	}

	merged := make(map[string]string, len(existing))
	for k, v := range existing {
		merged[k] = v
	}

	var inherited, skipped []string
	for k, v := range source {
		if _, exists := existing[k]; exists && !opts.Force {
			skipped = append(skipped, k)
			continue
		}
		merged[k] = v
		inherited = append(inherited, k)
	}

	if !opts.DryRun && len(inherited) > 0 {
		if err := client.Write(path, merged); err != nil {
			return Result{}, err
		}
	}

	return Result{
		Path:        path,
		Inherited:   inherited,
		Skipped:     skipped,
		DryRun:      opts.DryRun,
		InheritedAt: time.Now().UTC(),
	}, nil
}

func filterKeys(secrets map[string]string, keys []string) map[string]string {
	if len(keys) == 0 {
		return secrets
	}
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		if v, ok := secrets[k]; ok {
			out[k] = v
		}
	}
	return out
}
