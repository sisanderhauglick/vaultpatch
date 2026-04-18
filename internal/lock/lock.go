// Package lock provides advisory locking for Vault secret paths.
package lock

import (
	"errors"
	"fmt"
	"time"

	"github.com/your-org/vaultpatch/internal/vault"
)

// LockEntry represents a lock stored at a Vault path.
type LockEntry struct {
	Owner     string    `json:"owner"`
	LockedAt  time.Time `json:"locked_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Result holds the outcome of a lock or unlock operation.
type Result struct {
	Path    string
	Acquired bool
	Released bool
	DryRun  bool
	Entry   *LockEntry
}

// Lock acquires an advisory lock on the given Vault path.
func Lock(client *vault.Client, path, owner string, ttl time.Duration, dryRun bool) (Result, error) {
	if client == nil {
		return Result{}, errors.New("lock: client must not be nil")
	}
	if path == "" {
		return Result{}, errors.New("lock: path must not be empty")
	}
	if owner == "" {
		return Result{}, errors.New("lock: owner must not be empty")
	}

	now := time.Now().UTC()
	entry := &LockEntry{
		Owner:     owner,
		LockedAt:  now,
		ExpiresAt: now.Add(ttl),
	}

	if dryRun {
		return Result{Path: path, Acquired: true, DryRun: true, Entry: entry}, nil
	}

	lockPath := lockPath(path)
	existing, err := client.Read(lockPath)
	if err == nil && existing != nil {
		return Result{Path: path, Acquired: false, DryRun: false}, fmt.Errorf("lock: path %q is already locked", path)
	}

	data := map[string]interface{}{
		"owner":      entry.Owner,
		"locked_at":  entry.LockedAt.Format(time.RFC3339),
		"expires_at": entry.ExpiresAt.Format(time.RFC3339),
	}
	if err := client.Write(lockPath, data); err != nil {
		return Result{}, fmt.Errorf("lock: failed to write lock: %w", err)
	}

	return Result{Path: path, Acquired: true, DryRun: false, Entry: entry}, nil
}

// Unlock releases an advisory lock on the given Vault path.
func Unlock(client *vault.Client, path, owner string, dryRun bool) (Result, error) {
	if client == nil {
		return Result{}, errors.New("unlock: client must not be nil")
	}
	if path == "" {
		return Result{}, errors.New("unlock: path must not be empty")
	}

	if dryRun {
		return Result{Path: path, Released: true, DryRun: true}, nil
	}

	lockPath := lockPath(path)
	if err := client.Delete(lockPath); err != nil {
		return Result{}, fmt.Errorf("unlock: failed to delete lock: %w", err)
	}

	return Result{Path: path, Released: true, DryRun: false}, nil
}

func lockPath(path string) string {
	return path + "/.lock"
}
