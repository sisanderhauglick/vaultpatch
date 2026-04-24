// Package scope provides functionality to list and filter Vault secret paths
// matching a given prefix or glob pattern across one or more mounts.
package scope

import (
	"errors"
	"strings"
	"time"

	"github.com/your-org/vaultpatch/internal/vault"
)

// Result holds the outcome of a Scope operation for a single root path.
type Result struct {
	Root      string
	Paths     []string
	MatchedAt time.Time
	DryRun    bool
}

// Options controls how Scope behaves.
type Options struct {
	// Prefix filters paths to those starting with this value (case-insensitive).
	Prefix string
	// Keys restricts results to paths that contain at least one of these keys.
	Keys []string
	// DryRun returns results without performing any writes.
	DryRun bool
}

// Scope lists all secret paths under each root that match the given options.
func Scope(client *vault.Client, roots []string, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("scope: client must not be nil")
	}
	if len(roots) == 0 {
		return nil, errors.New("scope: at least one root path is required")
	}

	var results []Result
	for _, root := range roots {
		paths, err := client.List(root)
		if err != nil {
			return nil, err
		}

		matched := filterPaths(paths, opts.Prefix)

		if len(opts.Keys) > 0 {
			matched = filterByKeys(client, matched, opts.Keys)
		}

		results = append(results, Result{
			Root:      root,
			Paths:     matched,
			MatchedAt: time.Now().UTC(),
			DryRun:    opts.DryRun,
		})
	}
	return results, nil
}

// filterPaths returns only those paths that have the given prefix (case-insensitive).
// If prefix is empty all paths are returned.
func filterPaths(paths []string, prefix string) []string {
	if prefix == "" {
		return paths
	}
	lower := strings.ToLower(prefix)
	var out []string
	for _, p := range paths {
		if strings.HasPrefix(strings.ToLower(p), lower) {
			out = append(out, p)
		}
	}
	return out
}

// filterByKeys returns only those paths whose secret data contains at least one
// of the requested keys.
func filterByKeys(client *vault.Client, paths []string, keys []string) []string {
	var out []string
	for _, p := range paths {
		secrets, err := client.Read(p)
		if err != nil {
			continue
		}
		for _, k := range keys {
			if _, ok := secrets[k]; ok {
				out = append(out, p)
				break
			}
		}
	}
	return out
}
