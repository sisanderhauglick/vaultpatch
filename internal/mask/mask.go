// Package mask provides utilities for redacting sensitive secret values
// before display or export.
package mask

import "strings"

// Options controls masking behaviour.
type Options struct {
	// Keys is the list of key names to redact. If empty, all keys are masked.
	Keys []string
	// Placeholder replaces the original value. Defaults to "***".
	Placeholder string
}

// Result holds the masked secret map and metadata.
type Result struct {
	Path    string
	Masked  map[string]string
	Redacted int
}

// Mask redacts values in secrets according to opts.
func Mask(path string, secrets map[string]string, opts Options) Result {
	ph := opts.Placeholder
	if ph == "" {
		ph = "***"
	}

	out := make(map[string]string, len(secrets))
	redacted := 0

	for k, v := range secrets {
		if shouldMask(k, opts.Keys) {
			out[k] = ph
			redacted++
		} else {
			out[k] = v
		}
	}

	return Result{Path: path, Masked: out, Redacted: redacted}
}

func shouldMask(key string, keys []string) bool {
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
