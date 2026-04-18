// Package watch polls Vault paths for changes and emits diffs.
package watch

import (
	"context"
	"time"

	"github.com/your-org/vaultpatch/internal/diff"
	"github.com/your-org/vaultpatch/internal/vault"
)

// Event holds a detected change at a specific path.
type Event struct {
	Path    string
	Changes []diff.Change
	At      time.Time
}

// Options configures the watcher.
type Options struct {
	Paths    []string
	Interval time.Duration
}

// Watch polls the given paths at the specified interval, sending Events on the
// returned channel. The channel is closed when ctx is cancelled.
func Watch(ctx context.Context, client *vault.Client, opts Options) (<-chan Event, error) {
	if client == nil {
		return nil, fmt.Errorf("watch: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, fmt.Errorf("watch: at least one path is required")
	}
	if opts.Interval <= 0 {
		opts.Interval = 30 * time.Second
	}

	ch := make(chan Event, 8)

	go func() {
		defer close(ch)
		prev := make(map[string]vault.SecretMap)

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(opts.Interval):
				for _, path := range opts.Paths {
					current, err := vault.ReadSecrets(client, path)
					if err != nil {
						continue
					}
					changes := diff.Diff(prev[path], current)
					if len(changes) > 0 {
						ch <- Event{Path: path, Changes: changes, At: time.Now()}
					}
					prev[path] = current
				}
			}
		}
	}()

	return ch, nil
}
