// Package copy provides functionality to copy secrets between Vault paths.
package copy

import (
	"errors"
	"fmt"

	"github.com/your-org/vaultpatch/internal/vault"
)

// Result holds the outcome of a Copy operation.
type Result struct {
	SourcePath string
	DestPath   string
	Keys       int
	DryRun     bool
}

// Options configures a Copy operation.
type Options struct {
	SourcePath string
	DestPath   string
	Keys       []string // if empty, all keys are copied
	DryRun     bool
}

// Copy reads secrets from SourcePath and writes them to DestPath.
func Copy(client *vault.Client, opts Options) (Result, error) {
	if client == nil {
		return Result{}, errors.New("copy: client must not be nil")
	}
	if opts.SourcePath == "" || opts.DestPath == "" {
		return Result{}, errors.New("copy: source and dest paths must not be empty")
	}

	secrets, err := vault.ReadSecrets(client, opts.SourcePath)
	if err != nil {
		return Result{}, fmt.Errorf("copy: read source: %w", err)
	}

	filtered := filterKeys(secrets, opts.Keys)

	if !opts.DryRun {
		if err := vault.WriteSecrets(client, opts.DestPath, filtered); err != nil {
			return Result{}, fmt.Errorf("copy: write dest: %w", err)
		}
	}

	return Result{
		SourcePath: opts.SourcePath,
		DestPath:   opts.DestPath,
		Keys:       len(filtered),
		DryRun:     opts.DryRun,
	}, nil
}

func filterKeys(secrets vault.SecretMap, keys []string) vault.SecretMap {
	if len(keys) == 0 {
		return secrets
	}
	out := make(vault.SecretMap, len(keys))
	for _, k := range keys {
		if v, ok := secrets[k]; ok {
			out[k] = v
		}
	}
	return out
}
