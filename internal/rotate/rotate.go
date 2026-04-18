// Package rotate provides secret rotation helpers for vaultpatch.
package rotate

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/vaultpatch/vaultpatch/internal/vault"
)

// Result holds the outcome of a rotation operation.
type Result struct {
	Path    string
	Key     string
	DryRun  bool
	Rotated bool
}

// Options configures a Rotate call.
type Options struct {
	Path   string
	Keys   []string
	Length int
	DryRun bool
}

// Rotate generates new random values for the specified keys at the given path.
// If DryRun is true, no writes are performed.
func Rotate(client *vault.Client, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("rotate: client must not be nil")
	}
	if opts.Path == "" {
		return nil, errors.New("rotate: path must not be empty")
	}
	if len(opts.Keys) == 0 {
		return nil, errors.New("rotate: at least one key must be specified")
	}
	if opts.Length <= 0 {
		opts.Length = 32
	}

	current, err := vault.ReadSecrets(client, opts.Path)
	if err != nil {
		return nil, fmt.Errorf("rotate: read %s: %w", opts.Path, err)
	}

	updated := make(map[string]string, len(current))
	for k, v := range current {
		updated[k] = v
	}

	var results []Result
	for _, key := range opts.Keys {
		newVal, err := generateSecret(opts.Length)
		if err != nil {
			return nil, fmt.Errorf("rotate: generate value for key %q: %w", key, err)
		}
		if !opts.DryRun {
			updated[key] = newVal
		}
		results = append(results, Result{
			Path:    opts.Path,
			Key:     key,
			DryRun:  opts.DryRun,
			Rotated: !opts.DryRun,
		})
	}

	if !opts.DryRun {
		if err := vault.WriteSecrets(client, opts.Path, updated); err != nil {
			return nil, fmt.Errorf("rotate: write %s: %w", opts.Path, err)
		}
	}

	return results, nil
}

func generateSecret(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b)[:length], nil
}
