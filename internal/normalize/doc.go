// Package normalize provides functionality for standardising Vault secret
// values and key names in-place across one or more KV paths.
//
// Supported transformations include trimming whitespace, lowercasing key names,
// and uppercasing values. All operations support a dry-run mode that returns
// the projected changes without writing to Vault.
package normalize
