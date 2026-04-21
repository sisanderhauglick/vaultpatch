// Package diff2 provides structured two-way diffing between two Vault secret paths.
package diff2

import (
	"errors"
	"fmt"
	"time"

	"github.com/your-org/vaultpatch/internal/vault"
)

// ChangeKind describes the type of difference found.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
)

// Change represents a single key-level difference between two paths.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// Result holds the outcome of a two-path diff operation.
type Result struct {
	Source  string
	Dest    string
	Changes []Change
	DiffedAt time.Time
	DryRun  bool
}

// HasChanges returns true if any differences were found.
func (r Result) HasChanges() bool {
	return len(r.Changes) > 0
}

// Options configures a Diff2 run.
type Options struct {
	Source string
	Dest   string
	Keys   []string
	DryRun bool
}

// Diff2 reads secrets from two Vault paths and returns a structured diff result.
func Diff2(client *vault.Client, opts Options) (Result, error) {
	if client == nil {
		return Result{}, errors.New("diff2: client must not be nil")
	}
	if opts.Source == "" {
		return Result{}, errors.New("diff2: source path must not be empty")
	}
	if opts.Dest == "" {
		return Result{}, errors.New("diff2: destination path must not be empty")
	}

	src, err := vault.ReadSecrets(client, opts.Source)
	if err != nil {
		return Result{}, fmt.Errorf("diff2: reading source: %w", err)
	}
	dst, err := vault.ReadSecrets(client, opts.Dest)
	if err != nil {
		return Result{}, fmt.Errorf("diff2: reading dest: %w", err)
	}

	changes := computeChanges(src, dst, opts.Keys)

	return Result{
		Source:   opts.Source,
		Dest:     opts.Dest,
		Changes:  changes,
		DiffedAt: time.Now().UTC(),
		DryRun:   opts.DryRun,
	}, nil
}

func computeChanges(src, dst map[string]string, keys []string) []Change {
	filter := toSet(keys)
	want := func(k string) bool {
		if len(filter) == 0 {
			return true
		}
		_, ok := filter[k]
		return ok
	}

	var changes []Change
	for k, sv := range src {
		if !want(k) {
			continue
		}
		if dv, ok := dst[k]; !ok {
			changes = append(changes, Change{Key: k, Kind: Removed, OldValue: sv})
		} else if sv != dv {
			changes = append(changes, Change{Key: k, Kind: Modified, OldValue: sv, NewValue: dv})
		}
	}
	for k, dv := range dst {
		if !want(k) {
			continue
		}
		if _, ok := src[k]; !ok {
			changes = append(changes, Change{Key: k, Kind: Added, NewValue: dv})
		}
	}
	return changes
}

func toSet(keys []string) map[string]struct{} {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return m
}
