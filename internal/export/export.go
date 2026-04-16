package export

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

// Format represents the output format for exported secrets.
type Format string

const (
	FormatJSON Format = "json"
	FormatYAML Format = "yaml"
	FormatEnv  Format = "env"
)

// Export writes the secret map to w in the specified format.
func Export(secrets map[string]string, format Format, w io.Writer) error {
	switch format {
	case FormatJSON:
		return exportJSON(secrets, w)
	case FormatYAML:
		return exportYAML(secrets, w)
	case FormatEnv:
		return exportEnv(secrets, w)
	default:
		return fmt.Errorf("unsupported format: %q", format)
	}
}

func exportJSON(secrets map[string]string, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(secrets)
}

func exportYAML(secrets map[string]string, w io.Writer) error {
	return yaml.NewEncoder(w).Encode(secrets)
}

func exportEnv(secrets map[string]string, w io.Writer) error {
	for k, v := range secrets {
		key := strings.ToUpper(k)
		_, err := fmt.Fprintf(w, "%s=%q\n", key, v)
		if err != nil {
			return err
		}
	}
	return nil
}
