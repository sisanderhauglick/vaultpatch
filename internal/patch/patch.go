// Package patch applies a set of key-level changes to one or more Vault paths.
package patch

import (
	"errors"
	"fmt"
	"time"

	"github.com/youorg/vaultpatch/internal/vault"
)

// Op represents a single patch operation.
type Op struct {
	Key    string
	Value  string
	Delete bool
}

// Result holds the outcome of a Patch call.
type Result struct {
	Path      string
	Applied   []Op
	Skipped   []Op
	DryRun    bool
	PatchedAt time.Time
}

// Patch applies ops to each path. If dryRun is true, no writes occur.
func Patch(client *vault.Client, paths []string, ops []Op, dryRun bool) ([]Result, error) {
	if client == nil {
		return nil, errors.New("patch: client must not be nil")
	}
	if len(paths) == 0 {
		return nil, errors.New("patch: at least one path is required")
	}
	if len(ops) == 0 {
		return nil, errors.New("patch: at least one op is required")
	}

	var results []Result
	for _, p := range paths {
		secrets, err := vault.ReadSecrets(client, p)
		if err != nil {
			return nil, fmt.Errorf("patch: read %q: %w", p, err)
		}

		applied := []Op{}
		for _, op := range ops {
			if op.Delete {
				delete(secrets, op.Key)
			} else {
				secrets[op.Key] = op.Value
			}
			applied = append(applied, op)
		}

		if !dryRun {
			if err := vault.WriteSecrets(client, p, secrets); err != nil {
				return nil, fmt.Errorf("patch: write %q: %w", p, err)
			}
		}

		results = append(results, Result{
			Path:      p,
			Applied:   applied,
			DryRun:    dryRun,
			PatchedAt: time.Now().UTC(),
		})
	}
	return results, nil
}
