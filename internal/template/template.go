// Package template provides functionality for rendering Vault secret paths
// and values using Go text/template expressions.
package template

import (
	"bytes"
	"fmt"
	"strings"
	text "text/template"
)

// Result holds the output of a template rendering operation.
type Result struct {
	Path     string
	Rendered map[string]string
}

// Options controls how template rendering behaves.
type Options struct {
	// Vars are the key/value pairs injected into the template context.
	Vars map[string]string
	// Strict causes rendering to fail on missing variables when true.
	Strict bool
}

// Render applies Go template expressions to each value in secrets using the
// provided variables. Path itself is also rendered so callers can use dynamic
// path segments.
func Render(path string, secrets map[string]string, opts Options) (Result, error) {
	if path == "" {
		return Result{}, fmt.Errorf("template: path must not be empty")
	}
	if opts.Vars == nil {
		opts.Vars = map[string]string{}
	}

	renderedPath, err := renderString(path, opts)
	if err != nil {
		return Result{}, fmt.Errorf("template: path render failed: %w", err)
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		renderedVal, err := renderString(v, opts)
		if err != nil {
			return Result{}, fmt.Errorf("template: key %q render failed: %w", k, err)
		}
		out[k] = renderedVal
	}

	return Result{Path: renderedPath, Rendered: out}, nil
}

func renderString(s string, opts Options) (string, error) {
	if !strings.Contains(s, "{{") {
		return s, nil
	}

	option := "missingkey=zero"
	if opts.Strict {
		option = "missingkey=error"
	}

	tmpl, err := text.New("").Option(option).Parse(s)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, opts.Vars); err != nil {
		return "", err
	}
	return buf.String(), nil
}
