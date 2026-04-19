// Package archive provides functionality to archive (soft-delete) secrets
// at given Vault paths by writing metadata before removal.
package archive

import (
	"errors"
	"fmt"
	"time"

	"github.com/your-org/vaultpatch/internal/vault"
)

// Result holds the outcome of an archive operation.
type Result struct {
	Path     string
	Archived bool
	DryRun   bool
	ArchivedAt time.Time
}

// Options configures an Archive call.
type Options struct {
	Paths      []string
	ArchiveDest string // destination prefix to store archived copies, e.g. "secret/archive"
	DryRun     bool
}

// Archive copies secrets to an archive path then deletes the originals.
func Archive(client *vault.Client, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("archive: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("archive: at least one path is required")
	}
	if opts.ArchiveDest == "" {
		return nil, errors.New("archive: archive destination must not be empty")
	}

	now := time.Now().UTC()
	var results []Result

	for _, path := range opts.Paths {
		data, err := vault.ReadSecrets(client, path)
		if err != nil {
			return nil, fmt.Errorf("archive: read %q: %w", path, err)
		}

		destPath := fmt.Sprintf("%s/%s", opts.ArchiveDest, path)

		if !opts.DryRun {
			if err := vault.WriteSecrets(client, destPath, data); err != nil {
				return nil, fmt.Errorf("archive: write to %q: %w", destPath, err)
			}
			if err := vault.WriteSecrets(client, path, nil); err != nil {
				return nil, fmt.Errorf("archive: delete %q: %w", path, err)
			}
		}

		results = append(results, Result{
			Path:       path,
			Archived:   true,
			DryRun:     opts.DryRun,
			ArchivedAt: now,
		})
	}

	return results, nil
}
