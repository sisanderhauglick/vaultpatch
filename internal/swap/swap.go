// Package swap provides functionality to atomically swap secret values
// between two Vault paths, optionally performing a dry run.
package swap

import (
	"errors"
	"fmt"
	"time"

	"github.com/your-org/vaultpatch/internal/vault"
)

// Result holds the outcome of a Swap operation.
type Result struct {
	Source      string
	Destination string
	SwappedKeys []string
	SwappedAt   time.Time
	DryRun      bool
}

// Options controls Swap behaviour.
type Options struct {
	Keys   []string // if empty, all keys are swapped
	DryRun bool
}

// Swap exchanges secret key-values between src and dst paths.
// When DryRun is true no writes are performed.
func Swap(client *vault.Client, src, dst string, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("swap: client must not be nil")
	}
	if src == "" {
		return nil, errors.New("swap: source path must not be empty")
	}
	if dst == "" {
		return nil, errors.New("swap: destination path must not be empty")
	}

	srcData, err := client.Read(src)
	if err != nil {
		return nil, fmt.Errorf("swap: read source %q: %w", src, err)
	}
	dstData, err := client.Read(dst)
	if err != nil {
		return nil, fmt.Errorf("swap: read destination %q: %w", dst, err)
	}

	keys := opts.Keys
	if len(keys) == 0 {
		keys = mergedKeys(srcData, dstData)
	}

	newSrc := copyMap(srcData)
	newDst := copyMap(dstData)
	swapped := make([]string, 0, len(keys))

	for _, k := range keys {
		newSrc[k] = dstData[k]
		newDst[k] = srcData[k]
		swapped = append(swapped, k)
	}

	if !opts.DryRun {
		if err := client.Write(src, newSrc); err != nil {
			return nil, fmt.Errorf("swap: write source %q: %w", src, err)
		}
		if err := client.Write(dst, newDst); err != nil {
			return nil, fmt.Errorf("swap: write destination %q: %w", dst, err)
		}
	}

	return []Result{{
		Source:      src,
		Destination: dst,
		SwappedKeys: swapped,
		SwappedAt:   time.Now().UTC(),
		DryRun:      opts.DryRun,
	}}, nil
}

func mergedKeys(a, b map[string]interface{}) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out
}

func copyMap(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
