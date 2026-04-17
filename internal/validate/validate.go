package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule for a secret key or value.
type Rule struct {
	Key     string
	Pattern string
	Required bool
}

// Violation describes a single validation failure.
type Violation struct {
	Path    string
	Key     string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("%s/%s: %s", v.Path, v.Key, v.Message)
}

// Validate checks secrets at the given path against a set of rules.
// It returns a list of violations (empty means valid).
func Validate(path string, secrets map[string]string, rules []Rule) []Violation {
	var violations []Violation

	for _, rule := range rules {
		val, exists := secrets[rule.Key]

		if rule.Required && !exists {
			violations = append(violations, Violation{
				Path:    path,
				Key:     rule.Key,
				Message: "required key is missing",
			})
			continue
		}

		if !exists {
			continue
		}

		if strings.TrimSpace(val) == "" {
			violations = append(violations, Violation{
				Path:    path,
				Key:     rule.Key,
				Message: "value must not be empty or whitespace",
			})
			continue
		}

		if rule.Pattern != "" {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				violations = append(violations, Violation{
					Path:    path,
					Key:     rule.Key,
					Message: fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err),
				})
				continue
			}
			if !re.MatchString(val) {
				violations = append(violations, Violation{
					Path:    path,
					Key:     rule.Key,
					Message: fmt.Sprintf("value does not match pattern %q", rule.Pattern),
				})
			}
		}
	}

	return violations
}
