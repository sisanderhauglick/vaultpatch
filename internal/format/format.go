package format

import (
	"fmt"
	"sort"
	"strings"
)

// Style represents a display style for secret output.
type Style string

const (
	StyleTable Style = "table"
	StyleList  Style = "list"
	StyleCSV   Style = "csv"
)

// ParseStyle parses a string into a Style, returning an error for unknown values.
func ParseStyle(s string) (Style, error) {
	switch strings.ToLower(s) {
	case string(StyleTable):
		return StyleTable, nil
	case string(StyleList):
		return StyleList, nil
	case string(StyleCSV):
		return StyleCSV, nil
	default:
		return "", fmt.Errorf("unsupported style %q: choose from %s", s, strings.Join(SupportedStyles(), ", "))
	}
}

// SupportedStyles returns all valid style names.
func SupportedStyles() []string {
	return []string{string(StyleTable), string(StyleList), string(StyleCSV)}
}

// Render formats a map of secrets according to the given Style.
func Render(secrets map[string]string, style Style) string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	switch style {
	case StyleTable:
		sb.WriteString(fmt.Sprintf("%-30s %s\n", "KEY", "VALUE"))
		sb.WriteString(strings.Repeat("-", 50) + "\n")
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("%-30s %s\n", k, secrets[k]))
		}
	case StyleList:
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("%s=%s\n", k, secrets[k]))
		}
	case StyleCSV:
		sb.WriteString("key,value\n")
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("%s,%s\n", k, secrets[k]))
		}
	}
	return sb.String()
}
