package expire

import (
	"errors"
	"fmt"
	"time"

	"github.com/your-org/vaultpatch/internal/vault"
)

// Result holds the outcome of an expiry check or purge for a single path.
type Result struct {
	Path      string
	ExpiredAt time.Time
	Removed   bool
	DryRun    bool
}

// Options configures an Expire run.
type Options struct {
	Paths    []string
	Before   time.Time // purge secrets whose _expire_at metadata is before this time
	DryRun   bool
}

// Expire scans each path for an "_expire_at" key and removes secrets that have
// passed the deadline. In dry-run mode no writes are performed.
func Expire(client *vault.Client, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("expire: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("expire: at least one path is required")
	}
	if opts.Before.IsZero() {
		opts.Before = time.Now()
	}

	var results []Result
	for _, path := range opts.Paths {
		secrets, err := client.ReadSecrets(path)
		if err != nil {
			return nil, fmt.Errorf("expire: read %q: %w", path, err)
		}

		raw, ok := secrets["_expire_at"]
		if !ok {
			continue
		}

		expireAt, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", raw))
		if err != nil {
			return nil, fmt.Errorf("expire: parse _expire_at for %q: %w", path, err)
		}

		if expireAt.After(opts.Before) {
			continue
		}

		r := Result{Path: path, ExpiredAt: expireAt, DryRun: opts.DryRun}
		if !opts.DryRun {
			if err := client.DeleteSecrets(path); err != nil {
				return nil, fmt.Errorf("expire: delete %q: %w", path, err)
			}
			r.Removed = true
		}
		results = append(results, r)
	}
	return results, nil
}
