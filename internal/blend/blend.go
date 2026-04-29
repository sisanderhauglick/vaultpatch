// Package blend merges secrets from multiple source paths into a single
// destination path, with configurable conflict resolution strategies.
package blend

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
)

// Strategy controls how key conflicts are resolved when blending.
type Strategy string

const (
	// StrategyFirst keeps the value from the first source that defines the key.
	StrategyFirst Strategy = "first"
	// StrategyLast overwrites with the value from the last source that defines the key.
	StrategyLast Strategy = "last"
)

// Options configures a Blend operation.
type Options struct {
	Sources     []string
	Destination string
	Strategy    Strategy
	Keys        []string // if empty, all keys are included
	DryRun      bool
}

// Result describes the outcome of a Blend operation.
type Result struct {
	Destination string
	BlendedKeys []string
	SkippedKeys []string
	DryRun      bool
	BlendedAt   time.Time
}

// VaultClient is the subset of the Vault API used by Blend.
type VaultClient interface {
	Read(path string) (*api.Secret, error)
	Write(path string, data map[string]interface{}) (*api.Secret, error)
}

// Blend reads secrets from all source paths and writes the merged result to
// destination according to the chosen conflict resolution strategy.
func Blend(client VaultClient, opts Options) (Result, error) {
	if client == nil {
		return Result{}, errors.New("blend: client must not be nil")
	}
	if len(opts.Sources) == 0 {
		return Result{}, errors.New("blend: at least one source path is required")
	}
	if opts.Destination == "" {
		return Result{}, errors.New("blend: destination path must not be empty")
	}
	if opts.Strategy == "" {
		opts.Strategy = StrategyLast
	}

	merged := make(map[string]interface{})

	for _, src := range opts.Sources {
		secret, err := client.Read(src)
		if err != nil {
			return Result{}, fmt.Errorf("blend: read %q: %w", src, err)
		}
		if secret == nil || secret.Data == nil {
			continue
		}
		data := flatData(secret)
		for k, v := range data {
			if !shouldInclude(k, opts.Keys) {
				continue
			}
			_, exists := merged[k]
			if !exists || opts.Strategy == StrategyLast {
				merged[k] = v
			}
		}
	}

	blended := make([]string, 0, len(merged))
	for k := range merged {
		blended = append(blended, k)
	}

	res := Result{
		Destination: opts.Destination,
		BlendedKeys: blended,
		DryRun:      opts.DryRun,
		BlendedAt:   time.Now().UTC(),
	}

	if opts.DryRun {
		return res, nil
	}

	if _, err := client.Write(opts.Destination, merged); err != nil {
		return Result{}, fmt.Errorf("blend: write %q: %w", opts.Destination, err)
	}
	return res, nil
}

func flatData(s *api.Secret) map[string]interface{} {
	if kv2, ok := s.Data["data"]; ok {
		if m, ok := kv2.(map[string]interface{}); ok {
			return m
		}
	}
	return s.Data
}

func shouldInclude(key string, keys []string) bool {
	if len(keys) == 0 {
		return true
	}
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}
