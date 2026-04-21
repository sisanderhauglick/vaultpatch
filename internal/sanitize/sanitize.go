// Package sanitize provides utilities for cleaning and normalising Vault secret
// values before they are written back to a path (e.g. trimming whitespace,
// removing empty keys, or lower-casing key names).
package sanitize

import (
	"errors"
	"strings"
	"time"

	"github.com/your-org/vaultpatch/internal/vault"
)

// Options controls which sanitisation passes are applied.
type Options struct {
	TrimSpace    bool
	RemoveEmpty  bool
	LowercaseKeys bool
	Keys         []string // if non-empty, only these keys are processed
	DryRun       bool
}

// Result holds the outcome for a single Vault path.
type Result struct {
	Path        string
	Changed     map[string]string // key → new value (only mutated keys)
	Removed     []string
	DryRun      bool
	SanitizedAt time.Time
}

// Sanitize reads each path, applies the requested transformations, and writes
// the result back unless DryRun is true.
func Sanitize(client *vault.Client, paths []string, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("sanitize: client must not be nil")
	}
	if len(paths) == 0 {
		return nil, errors.New("sanitize: at least one path is required")
	}

	var results []Result
	for _, p := range paths {
		secrets, err := vault.ReadSecrets(client, p)
		if err != nil {
			return nil, err
		}

		res := Result{
			Path:        p,
			Changed:     make(map[string]string),
			DryRun:      opts.DryRun,
			SanitizedAt: time.Now().UTC(),
		}

		output := make(map[string]interface{})
		for k, v := range secrets {
			if !shouldProcess(k, opts.Keys) {
				output[k] = v
				continue
			}

			newKey := k
			if opts.LowercaseKeys {
				newKey = strings.ToLower(k)
			}

			strVal, _ := v.(string)
			if opts.TrimSpace {
				strVal = strings.TrimSpace(strVal)
			}
			if opts.RemoveEmpty && strVal == "" {
				res.Removed = append(res.Removed, k)
				continue
			}
			if newKey != k || strVal != v {
				res.Changed[newKey] = strVal
			}
			output[newKey] = strVal
		}

		if !opts.DryRun && (len(res.Changed) > 0 || len(res.Removed) > 0) {
			if err := vault.WriteSecrets(client, p, output); err != nil {
				return nil, err
			}
		}
		results = append(results, res)
	}
	return results, nil
}

func shouldProcess(key string, keys []string) bool {
	if len(keys) == 0 {
		return true
	}
	for _, k := range keys {
		if strings.EqualFold(k, key) {
			return true
		}
	}
	return false
}
