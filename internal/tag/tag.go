// Package tag provides functionality to read and write metadata tags
// on Vault secret paths using a reserved "_tags" key.
package tag

import (
	"errors"
	"fmt"
	"strings"

	"github.com/your-org/vaultpatch/internal/vault"
)

// Result holds the outcome of a tag operation.
type Result struct {
	Path    string
	Tags    []string
	DryRun  bool
	Updated bool
}

// Tag adds or removes tags on the given Vault path.
// Tags are stored as a comma-separated string under the key "_tags".
func Tag(client *vault.Client, path string, add, remove []string, dryRun bool) (Result, error) {
	if client == nil {
		return Result{}, errors.New("tag: client must not be nil")
	}
	if path == "" {
		return Result{}, errors.New("tag: path must not be empty")
	}

	secrets, err := vault.ReadSecrets(client, path)
	if err != nil {
		return Result{}, fmt.Errorf("tag: read %s: %w", path, err)
	}

	existing := parseTags(secrets["_tags"])
	merged := applyChanges(existing, add, remove)

	result := Result{
		Path:   path,
		Tags:   merged,
		DryRun: dryRun,
	}

	if dryRun {
		return result, nil
	}

	secrets["_tags"] = strings.Join(merged, ",")
	if err := vault.WriteSecrets(client, path, secrets); err != nil {
		return Result{}, fmt.Errorf("tag: write %s: %w", path, err)
	}
	result.Updated = true
	return result, nil
}

func parseTags(raw string) []string {
	if raw == "" {
		return []string{}
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func applyChanges(existing, add, remove []string) []string {
	set := make(map[string]struct{})
	for _, t := range existing {
		set[t] = struct{}{}
	}
	for _, t := range add {
		set[t] = struct{}{}
	}
	for _, t := range remove {
		delete(set, t)
	}
	out := make([]string, 0, len(set))
	for t := range set {
		out = append(out, t)
	}
	return out
}
