package snapshot

import (
	"errors"
	"fmt"
	"time"

	"github.com/your-org/vaultpatch/internal/vault"
)

// Snapshot represents a point-in-time capture of secrets at multiple paths.
type Snapshot struct {
	CapturedAt time.Time
	Paths      map[string]map[string]string
}

// Result holds the outcome of a snapshot operation.
type Result struct {
	Snapshot *Snapshot
	Paths    []string
	DryRun   bool
}

// Take captures current secrets at the given paths.
func Take(client *vault.Client, paths []string, dryRun bool) (*Result, error) {
	if client == nil {
		return nil, errors.New("vault client must not be nil")
	}
	if len(paths) == 0 {
		return nil, errors.New("at least one path is required")
	}

	snap := &Snapshot{
		CapturedAt: time.Now().UTC(),
		Paths:      make(map[string]map[string]string),
	}

	for _, p := range paths {
		secrets, err := vault.ReadSecrets(client, p)
		if err != nil {
			return nil, fmt.Errorf("reading path %q: %w", p, err)
		}
		snap.Paths[p] = secrets
	}

	return &Result{
		Snapshot: snap,
		Paths:    paths,
		DryRun:   dryRun,
	}, nil
}

// Diff compares a snapshot against current live secrets and returns changed paths.
func Diff(client *vault.Client, snap *Snapshot) (map[string][]string, error) {
	if client == nil {
		return nil, errors.New("vault client must not be nil")
	}
	if snap == nil {
		return nil, errors.New("snapshot must not be nil")
	}

	changed := make(map[string][]string)
	for path, old := range snap.Paths {
		current, err := vault.ReadSecrets(client, path)
		if err != nil {
			return nil, fmt.Errorf("reading path %q: %w", path, err)
		}
		for k, v := range current {
			if old[k] != v {
				changed[path] = append(changed[path], k)
			}
		}
		for k := range old {
			if _, ok := current[k]; !ok {
				changed[path] = append(changed[path], k)
			}
		}
	}
	return changed, nil
}
