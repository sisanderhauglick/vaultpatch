// Package resolve provides functionality to resolve secret references
// across Vault paths, expanding placeholder values like "ref:path/to/secret#key"
// into their concrete values.
package resolve

import (
	"fmt"
	"strings"
	"time"

	"github.com/your-org/vaultpatch/internal/vault"
)

// Result holds the outcome of resolving references within a single path.
type Result struct {
	Path       string
	Resolved   map[string]string
	Unresolved []string
	DryRun     bool
	ResolvedAt time.Time
}

// Options configures the Resolve operation.
type Options struct {
	Paths  []string
	DryRun bool
}

// Resolve reads secrets at each path and expands any value matching the
// pattern "ref:<vault-path>#<key>" by fetching the referenced secret.
func Resolve(client *vault.Client, opts Options) ([]Result, error) {
	if client == nil {
		return nil, fmt.Errorf("resolve: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, fmt.Errorf("resolve: at least one path is required")
	}

	var results []Result

	for _, path := range opts.Paths {
		secrets, err := vault.ReadSecrets(client, path)
		if err != nil {
			return nil, fmt.Errorf("resolve: read %q: %w", path, err)
		}

		resolved := make(map[string]string)
		var unresolved []string

		for k, v := range secrets {
			str, ok := v.(string)
			if !ok {
				continue
			}
			if !strings.HasPrefix(str, "ref:") {
				continue
			}

			refPath, refKey, err := parseRef(str)
			if err != nil {
				unresolved = append(unresolved, k)
				continue
			}

			refSecrets, err := vault.ReadSecrets(client, refPath)
			if err != nil || refSecrets[refKey] == nil {
				unresolved = append(unresolved, k)
				continue
			}

			resolved[k] = fmt.Sprintf("%v", refSecrets[refKey])
		}

		if !opts.DryRun && len(resolved) > 0 {
			merged := make(map[string]interface{})
			for k, v := range secrets {
				merged[k] = v
			}
			for k, v := range resolved {
				merged[k] = v
			}
			if err := vault.WriteSecrets(client, path, merged); err != nil {
				return nil, fmt.Errorf("resolve: write %q: %w", path, err)
			}
		}

		results = append(results, Result{
			Path:       path,
			Resolved:   resolved,
			Unresolved: unresolved,
			DryRun:     opts.DryRun,
			ResolvedAt: time.Now().UTC(),
		})
	}

	return results, nil
}

// parseRef splits a reference string "ref:<path>#<key>" into its components.
func parseRef(ref string) (path, key string, err error) {
	without := strings.TrimPrefix(ref, "ref:")
	parts := strings.SplitN(without, "#", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid ref format: %q", ref)
	}
	return parts[0], parts[1], nil
}
