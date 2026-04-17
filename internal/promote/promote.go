package promote

import (
	"fmt"

	"github.com/vaultpatch/internal/diff"
	"github.com/vaultpatch/internal/vault"
)

// Options configures a promotion run.
type Options struct {
	SrcPath  string
	DstPath  string
	DryRun   bool
	Keys     []string // if non-empty, only promote these keys
}

// Result holds the outcome of a promotion.
type Result struct {
	Changes []diff.Change
	Applied int
	Skipped int
}

// Promote copies secrets from src path to dst path, returning a Result.
func Promote(client *vault.Client, opts Options) (*Result, error) {
	if client == nil {
		return nil, fmt.Errorf("promote: vault client is nil")
	}
	if opts.SrcPath == "" || opts.DstPath == "" {
		return nil, fmt.Errorf("promote: src and dst paths must not be empty")
	}

	src, err := vault.ReadSecrets(client, opts.SrcPath)
	if err != nil {
		return nil, fmt.Errorf("promote: read src: %w", err)
	}

	dst, err := vault.ReadSecrets(client, opts.DstPath)
	if err != nil {
		return nil, fmt.Errorf("promote: read dst: %w", err)
	}

	// Filter src to requested keys.
	filtered := filterKeys(src, opts.Keys)

	changes := diff.Diff(filtered, dst)

	result := &Result{Changes: changes}

	for _, c := range changes {
		if opts.DryRun {
			result.Skipped++
			continue
		}
		dst[c.Key] = c.NewValue
		result.Applied++
	}

	if !opts.DryRun && result.Applied > 0 {
		if err := vault.WriteSecrets(client, opts.DstPath, dst); err != nil {
			return nil, fmt.Errorf("promote: write dst: %w", err)
		}
	}

	return result, nil
}

func filterKeys(src map[string]string, keys []string) map[string]string {
	if len(keys) == 0 {
		return src
	}
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		if v, ok := src[k]; ok {
			out[k] = v
		}
	}
	return out
}
