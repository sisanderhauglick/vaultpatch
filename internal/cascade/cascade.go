// Package cascade propagates secret values from a source path to one or more
// destination paths, merging only the specified keys (or all keys when none
// are given). Each destination write is skipped in dry-run mode.
package cascade

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
)

// VaultClient is the subset of the Vault API used by Cascade.
type VaultClient interface {
	Read(path string) (*api.Secret, error)
	Write(path string, data map[string]interface{}) (*api.Secret, error)
}

// Result describes the outcome of cascading a source to a single destination.
type Result struct {
	Source      string
	Destination string
	KeysCopied  int
	DryRun      bool
	CascadedAt  time.Time
}

// Options controls how Cascade behaves.
type Options struct {
	Source       string
	Destinations []string
	Keys         []string // empty means all keys
	DryRun       bool
}

// Cascade reads secrets from Source and writes the selected keys to every
// Destination path, returning one Result per destination.
func Cascade(client VaultClient, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("cascade: client must not be nil")
	}
	if opts.Source == "" {
		return nil, errors.New("cascade: source path must not be empty")
	}
	if len(opts.Destinations) == 0 {
		return nil, errors.New("cascade: at least one destination path is required")
	}

	secret, err := client.Read(opts.Source)
	if err != nil {
		return nil, fmt.Errorf("cascade: read source %q: %w", opts.Source, err)
	}

	data := secretData(secret)
	selected := filterKeys(data, opts.Keys)

	var results []Result
	for _, dest := range opts.Destinations {
		r := Result{
			Source:      opts.Source,
			Destination: dest,
			KeysCopied:  len(selected),
			DryRun:      opts.DryRun,
			CascadedAt:  time.Now().UTC(),
		}
		if !opts.DryRun && len(selected) > 0 {
			if _, err := client.Write(dest, selected); err != nil {
				return nil, fmt.Errorf("cascade: write destination %q: %w", dest, err)
			}
		}
		results = append(results, r)
	}
	return results, nil
}

func secretData(s *api.Secret) map[string]interface{} {
	if s == nil || s.Data == nil {
		return map[string]interface{}{}
	}
	return s.Data
}

func filterKeys(data map[string]interface{}, keys []string) map[string]interface{} {
	if len(keys) == 0 {
		out := make(map[string]interface{}, len(data))
		for k, v := range data {
			out[k] = v
		}
		return out
	}
	out := make(map[string]interface{}, len(keys))
	for _, k := range keys {
		if v, ok := data[k]; ok {
			out[k] = v
		}
	}
	return out
}
