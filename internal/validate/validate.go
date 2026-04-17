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

// Summary returns a single error summarising all violations, or nil if there are none.
func Summary(violations []Violation) error {
	if len(violations) == 0 {
		return nil
	}
	msgs := make([]string, len(violations))
	for i, v := range violations {
		msgs[i] = v.Error()
	}
	return fmt.Errorf("%d validation violation(s):\n%s", len(violations), strings.Join(msgs, "\n"))
}
