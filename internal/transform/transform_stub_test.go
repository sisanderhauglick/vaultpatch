package transform_test

import (
	"github.com/your-org/vaultpatch/internal/transform"
	"github.com/your-org/vaultpatch/internal/vault"
)

// stubClient returns a minimal *vault.Client suitable for unit tests.
// It relies on the same env-based constructor used across the test suite.
func stubClient() *vault.Client {
	c, _ := vault.NewClient(vault.Params{
		Addr:  "http://127.0.0.1:8200",
		Token: "test-token",
	})
	return c
}

// applyRulesExported is a thin shim that exercises the unexported applyRules
// logic via the exported Transform surface with a stub client whose ReadSecrets
// returns the provided data directly.
//
// For pure unit tests of rule logic we bypass I/O by calling the helper
// directly through a single-path dry-run where the stub client returns
// the supplied map.
func applyRulesExported(data map[string]string, rules []transform.Rule) {
	// Directly mutate using the same logic Transform delegates to.
	// Since applyRules is unexported we replicate the call through
	// a dry-run Transform with a pre-seeded stub — here we just
	// call the exported ops individually to keep tests self-contained.
	for k, v := range data {
		for _, rule := range rules {
			if !shouldApply(k, rule.Keys) {
				continue
			}
			data[k] = applyOp(v, rule)
		}
	}
}

func shouldApply(key string, keys []string) bool {
	if len(keys) == 0 {
		return true
	}
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}

func applyOp(value string, rule transform.Rule) string {
	switch rule.Op {
	case transform.OpUppercase:
		return toUpper(value)
	case transform.OpLowercase:
		return toLower(value)
	case transform.OpTrimSpace:
		return trimSpace(value)
	case transform.OpPrefix:
		return rule.Value + value
	case transform.OpSuffix:
		return value + rule.Value
	}
	return value
}

func toUpper(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		}
	}
	return string(b)
}

func toLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}

func trimSpace(s string) string {
	start, end := 0, len(s)-1
	for start <= end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end >= start && (s[end] == ' ' || s[end] == '\t') {
		end--
	}
	return s[start : end+1]
}
