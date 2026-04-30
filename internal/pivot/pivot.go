// Package pivot reorganises secrets from multiple source paths into a
// single destination path, keying each value by a configurable prefix
// derived from the source path segment.
package pivot

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
)

// VaultClient is the subset of the Vault API used by Pivot.
type VaultClient interface {
	Read(path string) (map[string]interface{}, error)
	Write(path string, data map[string]interface{}) error
}

// Options controls Pivot behaviour.
type Options struct {
	// Sources is the ordered list of KV paths to read from.
	Sources []string
	// Destination is the KV path to write the merged result to.
	Destination string
	// KeyPrefix, when non-empty, is prepended to each key as "<prefix>_<key>".
	// When empty the last segment of the source path is used.
	KeyPrefix string
	// DryRun skips the final write when true.
	DryRun bool
}

// Result describes the outcome of a Pivot operation.
type Result struct {
	Destination string
	MergedKeys  []string
	PivotedAt   time.Time
	DryRun      bool
}

// Pivot reads secrets from each source path and writes them to destination,
// namespacing every key with a prefix derived from the source path.
func Pivot(client VaultClient, opts Options) (*Result, error) {
	if client == nil {
		return nil, errors.New("pivot: client must not be nil")
	}
	if len(opts.Sources) == 0 {
		return nil, errors.New("pivot: at least one source path is required")
	}
	if opts.Destination == "" {
		return nil, errors.New("pivot: destination path must not be empty")
	}

	merged := make(map[string]interface{})
	var mergedKeys []string

	for _, src := range opts.Sources {
		data, err := client.Read(src)
		if err != nil {
			return nil, fmt.Errorf("pivot: read %q: %w", src, err)
		}
		prefix := opts.KeyPrefix
		if prefix == "" {
			parts := strings.Split(strings.TrimRight(src, "/"), "/")
			prefix = parts[len(parts)-1]
		}
		for k, v := range data {
			newKey := prefix + "_" + k
			merged[newKey] = v
			mergedKeys = append(mergedKeys, newKey)
		}
	}

	if !opts.DryRun {
		if err := client.Write(opts.Destination, merged); err != nil {
			return nil, fmt.Errorf("pivot: write %q: %w", opts.Destination, err)
		}
	}

	return &Result{
		Destination: opts.Destination,
		MergedKeys:  mergedKeys,
		PivotedAt:   time.Now().UTC(),
		DryRun:      opts.DryRun,
	}, nil
}

// ensure api.Client is not accidentally imported at compile time.
var _ *api.Client
