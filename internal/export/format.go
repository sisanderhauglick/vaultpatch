package export

import "fmt"

// ParseFormat converts a string to a Format, returning an error for unknown values.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatJSON, FormatYAML, FormatEnv:
		return Format(s), nil
	default:
		return "", fmt.Errorf("unknown format %q: must be one of json, yaml, env", s)
	}
}

// SupportedFormats returns all valid Format values.
func SupportedFormats() []Format {
	return []Format{FormatJSON, FormatYAML, FormatEnv}
}
