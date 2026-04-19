package redact

import "strings"

// Result holds the redacted secrets for a single path.
type Result struct {
	Path     string
	Original map[string]string
	Redacted map[string]string
}

// Options configures redaction behaviour.
type Options struct {
	Keys        []string // if empty, all keys are redacted
	Replacement string   // defaults to "[REDACTED]"
}

// Redact replaces secret values with a placeholder.
// If opts.Keys is non-empty only those keys are redacted; others are passed through.
func Redact(secrets map[string]map[string]string, opts Options) []Result {
	if opts.Replacement == "" {
		opts.Replacement = "[REDACTED]"
	}

	var results []Result
	for path, kv := range secrets {
		original := make(map[string]string, len(kv))
		redacted := make(map[string]string, len(kv))
		for k, v := range kv {
			original[k] = v
			if shouldRedact(k, opts.Keys) {
				redacted[k] = opts.Replacement
			} else {
				redacted[k] = v
			}
		}
		results = append(results, Result{
			Path:     path,
			Original: original,
			Redacted: redacted,
		})
	}
	return results
}

func shouldRedact(key string, keys []string) bool {
	if len(keys) == 0 {
		return true
	}
	for _, k := range keys {
		if strings.EqualFold(k, key) {
			return true
		}
	}
	return false
}
