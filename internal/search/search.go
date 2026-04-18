package search

import (
	"fmt"
	"strings"

	"github.com/vaultpatch/internal/vault"
)

// Result holds a single match found during a search.
type Result struct {
	Path  string
	Key   string
	Value string
}

// Options controls search behaviour.
type Options struct {
	Paths       []string
	Query       string
	KeysOnly    bool
	CaseSensitive bool
}

// Search scans the given Vault paths for secrets whose keys or values contain
// the query string and returns all matches.
func Search(client *vault.Client, opts Options) ([]Result, error) {
	if client == nil {
		return nil, fmt.Errorf("search: vault client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, fmt.Errorf("search: at least one path is required")
	}
	if opts.Query == "" {
		return nil, fmt.Errorf("search: query must not be empty")
	}

	query := opts.Query
	if !opts.CaseSensitive {
		query = strings.ToLower(query)
	}

	var results []Result
	for _, path := range opts.Paths {
		secrets, err := vault.ReadSecrets(client, path)
		if err != nil {
			return nil, fmt.Errorf("search: reading %s: %w", path, err)
		}
		for k, v := range secrets {
			ck := k
			cv := v
			if !opts.CaseSensitive {
				ck = strings.ToLower(k)
				cv = strings.ToLower(v)
			}
			keyMatch := strings.Contains(ck, query)
			valMatch := !opts.KeysOnly && strings.Contains(cv, query)
			if keyMatch || valMatch {
				results = append(results, Result{Path: path, Key: k, Value: v})
			}
		}
	}
	return results, nil
}
