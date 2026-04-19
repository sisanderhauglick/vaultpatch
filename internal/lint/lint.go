package lint

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Rule defines a single lint rule applied to secret keys or values.
type Rule struct {
	Name    string
	Message string
	Check   func(key, value string) bool
}

// Violation describes a rule breach at a specific path and key.
type Violation struct {
	Path    string
	Key     string
	Rule    string
	Message string
}

// Result holds the outcome of a lint run.
type Result struct {
	Path       string
	Violations []Violation
	LintedAt   time.Time
}

// Options controls lint behaviour.
type Options struct {
	Paths []string
	Rules []Rule
}

// VaultClient is the subset of vault operations lint requires.
type VaultClient interface {
	Read(path string) (map[string]interface{}, error)
}

// DefaultRules returns the built-in lint rules.
func DefaultRules() []Rule {
	return []Rule{
		{
			Name:    "no-empty-value",
			Message: "secret value must not be empty",
			Check:   func(_, v string) bool { return strings.TrimSpace(v) == "" },
		},
		{
			Name:    "no-whitespace-key",
			Message: "secret key must not contain whitespace",
			Check:   func(k, _ string) bool { return regexp.MustCompile(`\s`).MatchString(k) },
		},
		{
			Name:    "no-uppercase-key",
			Message: "secret key should be lowercase",
			Check:   func(k, _ string) bool { return k != strings.ToLower(k) },
		},
	}
}

// Lint reads secrets at each path and applies rules, returning per-path results.
func Lint(client VaultClient, opts Options) ([]Result, error) {
	if client == nil {
		return nil, errors.New("lint: client must not be nil")
	}
	if len(opts.Paths) == 0 {
		return nil, errors.New("lint: at least one path is required")
	}
	if len(opts.Rules) == 0 {
		opts.Rules = DefaultRules()
	}

	var results []Result
	for _, path := range opts.Paths {
		data, err := client.Read(path)
		if err != nil {
			return nil, fmt.Errorf("lint: read %q: %w", path, err)
		}
		r := Result{Path: path, LintedAt: time.Now()}
		for k, v := range data {
			val := fmt.Sprintf("%v", v)
			for _, rule := range opts.Rules {
				if rule.Check(k, val) {
					r.Violations = append(r.Violations, Violation{
						Path: path, Key: k,
						Rule: rule.Name, Message: rule.Message,
					})
				}
			}
		}
		results = append(results, r)
	}
	return results, nil
}
