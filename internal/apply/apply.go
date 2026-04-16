package apply

import (
	"fmt"

	"github.com/yourusername/vaultpatch/internal/diff"
	"github.com/yourusername/vaultpatch/internal/vault"
)

// Result holds the outcome of applying a single diff entry.
type Result struct {
	Key     string
	Action  string
	Success bool
	Err     error
}

// Options controls apply behaviour.
type Options struct {
	DryRun bool
}

// Apply writes diff changes to Vault at the given mount/path.
// Returns one Result per change entry.
func Apply(client *vault.Client, mount, path string, changes []diff.Entry, opts Options) []Result {
	results := make([]Result, 0, len(changes))

	// Build the full secret map by reading current state first.
	current, err := vault.ReadSecrets(client, mount, path)
	if err != nil {
		// Return a single error result covering all entries.
		return []Result{{Key: "*", Action: "read", Success: false, Err: err}}
	}

	for _, entry := range changes {
		r := Result{Key: entry.Key, Action: string(entry.Type)}

		switch entry.Type {
		case diff.Added, diff.Modified:
			current[entry.Key] = entry.NewValue
		case diff.Removed:
			delete(current, entry.Key)
		default:
			r.Err = fmt.Errorf("unknown diff type: %s", entry.Type)
			results = append(results, r)
			continue
		}

		if opts.DryRun {
			r.Success = true
			results = append(results, r)
			continue
		}

		if err := vault.WriteSecrets(client, mount, path, current); err != nil {
			r.Err = err
		} else {
			r.Success = true
		}
		results = append(results, r)
	}

	return results
}
