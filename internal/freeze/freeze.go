// Package freeze provides functionality to freeze and unfreeze Vault secret
// paths, preventing any writes until explicitly released.
package freeze

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
)

const freezeKey = "__frozen__"

// Client is the subset of the Vault API client used by this package.
type Client interface {
	Read(path string) (*api.Secret, error)
	Write(path string, data map[string]interface{}) (*api.Secret, error)
}

// Result holds the outcome of a Freeze or Unfreeze operation on a single path.
type Result struct {
	Path       string
	Frozen     bool
	FrozenAt   time.Time
	FrozenBy   string
	DryRun     bool
}

// Options configures a Freeze or Unfreeze call.
type Options struct {
	Paths  []string
	Owner  string
	DryRun bool
}

// Freeze marks each path as frozen by writing a metadata key.
func Freeze(client Client, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("freeze: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("freeze: at least one path is required")
	}
	if opts.Owner == "" {
		return nil, errors.New("freeze: owner must not be empty")
	}

	now := time.Now().UTC()
	results := make([]Result, 0, len(opts.Paths))

	for _, p := range opts.Paths {
		r := Result{
			Path:     p,
			Frozen:   true,
			FrozenAt: now,
			FrozenBy: opts.Owner,
			DryRun:   opts.DryRun,
		}
		if !opts.DryRun {
			existing, err := client.Read(p)
			if err != nil {
				return nil, fmt.Errorf("freeze: read %s: %w", p, err)
			}
			data := mergeData(existing, map[string]interface{}{
				freezeKey: fmt.Sprintf("%s@%s", opts.Owner, now.Format(time.RFC3339)),
			})
			if _, err := client.Write(p, data); err != nil {
				return nil, fmt.Errorf("freeze: write %s: %w", p, err)
			}
		}
		results = append(results, r)
	}
	return results, nil
}

// Unfreeze removes the freeze marker from each path.
func Unfreeze(client Client, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("freeze: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("freeze: at least one path is required")
	}

	results := make([]Result, 0, len(opts.Paths))
	for _, p := range opts.Paths {
		r := Result{Path: p, Frozen: false, DryRun: opts.DryRun}
		if !opts.DryRun {
			existing, err := client.Read(p)
			if err != nil {
				return nil, fmt.Errorf("freeze: read %s: %w", p, err)
			}
			data := mergeData(existing, nil)
			delete(data, freezeKey)
			if _, err := client.Write(p, data); err != nil {
				return nil, fmt.Errorf("freeze: write %s: %w", p, err)
			}
		}
		results = append(results, r)
	}
	return results, nil
}

// IsFrozen reports whether the secret at path carries the freeze marker.
func IsFrozen(client Client, path string) (bool, error) {
	if client == nil {
		return false, errors.New("freeze: client must not be nil")
	}
	s, err := client.Read(path)
	if err != nil {
		return false, fmt.Errorf("freeze: read %s: %w", path, err)
	}
	if s == nil || s.Data == nil {
		return false, nil
	}
	_, ok := s.Data[freezeKey]
	return ok, nil
}

func mergeData(existing *api.Secret, extra map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	if existing != nil {
		for k, v := range existing.Data {
			out[k] = v
		}
	}
	for k, v := range extra {
		out[k] = v
	}
	return out
}
