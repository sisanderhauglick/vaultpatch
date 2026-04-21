// Package flatten merges secrets from multiple Vault paths into a single
// destination path. It supports selective key filtering, conflict resolution
// via an overwrite flag, and dry-run mode to preview changes without writing.
package flatten
