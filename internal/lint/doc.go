// Package lint checks Vault secret paths against configurable rules,
// reporting violations such as empty values, whitespace in keys, or
// non-lowercase key names. Custom rules can be supplied alongside or
// instead of the built-in defaults.
package lint
