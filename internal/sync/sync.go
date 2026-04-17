// Package sync provides functionality to synchronize secrets between two Vault paths.
package sync

import (
	"fmt"

	"github.com/your-org/vaultpatch/internal/diff"
	"github.com/your-org/vaultpatch/internal/vault"
)

// Options configures a Sync operation.
type Options struct {
	SrcPath  string
	DstPath  string
	DryRun   bool
	Includes []string // if non-empty, only sync these keys
}

// Result summarises what happened during a Sync.
type Result struct {
	Changes []diff.Change
	Applied int
}

// Sync reads secrets from SrcPath, diffs them against DstPath, and writes
// any changes to DstPath (unless DryRun is true).
func Sync(client *vault.Client, opts Options) (*Result, error) {
	if client == nil {
		return nil, fmt.Errorf("sync: client must not be nil")
	}
	if opts.SrcPath == "" || opts.DstPath == "" {
		return nil, fmt.Errorf("sync: SrcPath and DstPath must not be empty")
	}

	src, err := vault.ReadSecrets(client, opts.SrcPath)
	if err != nil {
		return nil, fmt.Errorf("sync: read src: %w", err)
	}

	dst, err := vault.ReadSecrets(client, opts.DstPath)
	if err != nil {
		return nil, fmt.Errorf("sync: read dst: %w", err)
	}

	filtered := filterKeys(src, opts.Includes)
	changes := diff.Diff(filtered, dst)

	result := &Result{Changes: changes}
	if opts.DryRun || len(changes) == 0 {
		return result, nil
	}

	merged := applyChanges(dst, changes)
	if err := vault.WriteSecrets(client, opts.DstPath, merged); err != nil {
		return nil, fmt.Errorf("sync: write dst: %w", err)
	}
	result.Applied = len(changes)
	return result, nil
}

func filterKeys(src vault.SecretMap, includes []string) vault.SecretMap {
	if len(includes) == 0 {
		return src
	}
	out := make(vault.SecretMap, len(includes))
	for _, k := range includes {
		if v, ok := src[k]; ok {
			out[k] = v
		}
	}
	return out
}

func applyChanges(base vault.SecretMap, changes []diff.Change) vault.SecretMap {
	out := make(vault.SecretMap, len(base))
	for k, v := range base {
		out[k] = v
	}
	for _, c := range changes {
		switch c.Action {
		case diff.ActionAdded, diff.ActionModified:
			out[c.Key] = c.NewValue
		case diff.ActionRemoved:
			delete(out, c.Key)
		}
	}
	return out
}
