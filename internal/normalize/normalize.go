package normalize

import (
	"fmt"
	"strings"

	"github.com/your-org/vaultpatch/internal/vault"
)

// Options configures a Normalize operation.
type Options struct {
	Paths       []string
	Keys        []string // if empty, all keys are processed
	TrimSpace   bool
	LowercaseKeys bool
	UppercaseValues bool
	DryRun      bool
}

// Result holds the outcome of normalizing a single path.
type Result struct {
	Path     string
	Changes  map[string]string // key -> new value
	DryRun   bool
}

// Normalize reads secrets at each path, applies normalization rules, and
// writes the result back unless DryRun is set.
func Normalize(client *vault.Client, opts Options) ([]Result, error) {
	if client == nil {
		return nil, fmt.Errorf("normalize: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, fmt.Errorf("normalize: at least one path is required")
	}

	var results []Result

	for _, path := range opts.Paths {
		secrets, err := vault.ReadSecrets(client, path)
		if err != nil {
			return nil, fmt.Errorf("normalize: read %q: %w", path, err)
		}

		changes := applyRules(secrets, opts)
		if len(changes) == 0 {
			continue
		}

		if !opts.DryRun {
			merged := make(map[string]interface{}, len(secrets))
			for k, v := range secrets {
				merged[k] = v
			}
			for k, v := range changes {
				merged[k] = v
			}
			if err := vault.WriteSecrets(client, path, merged); err != nil {
				return nil, fmt.Errorf("normalize: write %q: %w", path, err)
			}
		}

		results = append(results, Result{
			Path:    path,
			Changes: changes,
			DryRun:  opts.DryRun,
		})
	}

	return results, nil
}

// applyRules returns a map of key->normalizedValue for any key whose value changed.
func applyRules(secrets map[string]interface{}, opts Options) map[string]string {
	changes := make(map[string]string)

	for rawKey, rawVal := range secrets {
		if !shouldProcess(rawKey, opts.Keys) {
			continue
		}

		key := rawKey
		if opts.LowercaseKeys {
			key = strings.ToLower(rawKey)
		}

		val := fmt.Sprintf("%v", rawVal)
		if opts.TrimSpace {
			val = strings.TrimSpace(val)
		}
		if opts.UppercaseValues {
			val = strings.ToUpper(val)
		}

		original := fmt.Sprintf("%v", rawVal)
		if key != rawKey || val != original {
			changes[key] = val
		}
	}

	return changes
}

// shouldProcess returns true when keys is empty or rawKey is in keys.
func shouldProcess(rawKey string, keys []string) bool {
	if len(keys) == 0 {
		return true
	}
	for _, k := range keys {
		if strings.EqualFold(k, rawKey) {
			return true
		}
	}
	return false
}
