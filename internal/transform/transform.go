// Package transform applies key/value transformations to Vault secrets.
package transform

import (
	"errors"
	"strings"
	"time"

	"github.com/your-org/vaultpatch/internal/vault"
)

// Op represents a transformation operation type.
type Op string

const (
	OpUppercase Op = "uppercase"
	OpLowercase Op = "lowercase"
	OpTrimSpace Op = "trimspace"
	OpPrefix    Op = "prefix"
	OpSuffix    Op = "suffix"
)

// Rule describes a single transformation to apply.
type Rule struct {
	Op    Op
	Keys  []string // empty means all keys
	Value string   // used by prefix/suffix ops
}

// Result holds the outcome of a Transform call for one path.
type Result struct {
	Path          string
	Transformed   map[string]string
	ChangedCount  int
	DryRun        bool
	TransformedAt time.Time
}

// Options controls Transform behaviour.
type Options struct {
	Paths  []string
	Rules  []Rule
	DryRun bool
}

// Transform reads secrets at each path, applies the given rules, and writes
// the results back unless DryRun is set.
func Transform(c *vault.Client, opts Options) ([]Result, error) {
	if c == nil {
		return nil, errors.New("transform: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("transform: at least one path is required")
	}
	if len(opts.Rules) == 0 {
		return nil, errors.New("transform: at least one rule is required")
	}

	var results []Result
	for _, path := range opts.Paths {
		secrets, err := vault.ReadSecrets(c, path)
		if err != nil {
			return nil, err
		}

		transformed := make(map[string]string, len(secrets))
		for k, v := range secrets {
			transformed[k] = v
		}

		changed := applyRules(transformed, opts.Rules)

		if !opts.DryRun {
			raw := make(map[string]interface{}, len(transformed))
			for k, v := range transformed {
				raw[k] = v
			}
			if err := vault.WriteSecrets(c, path, raw); err != nil {
				return nil, err
			}
		}

		results = append(results, Result{
			Path:          path,
			Transformed:   transformed,
			ChangedCount:  changed,
			DryRun:        opts.DryRun,
			TransformedAt: time.Now().UTC(),
		})
	}
	return results, nil
}

// applyRules mutates data in place and returns the number of changed values.
func applyRules(data map[string]string, rules []Rule) int {
	changed := 0
	for _, rule := range rules {
		for k, v := range data {
			if !shouldApply(k, rule.Keys) {
				continue
			}
			newVal := applyOp(v, rule)
			if newVal != v {
				data[k] = newVal
				changed++
			}
		}
	}
	return changed
}

func shouldApply(key string, keys []string) bool {
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

func applyOp(value string, rule Rule) string {
	switch rule.Op {
	case OpUppercase:
		return strings.ToUpper(value)
	case OpLowercase:
		return strings.ToLower(value)
	case OpTrimSpace:
		return strings.TrimSpace(value)
	case OpPrefix:
		return rule.Value + value
	case OpSuffix:
		return value + rule.Value
	}
	return value
}
