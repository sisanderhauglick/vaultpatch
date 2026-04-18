// Package clone implements deep-copy of Vault secret paths.
//
// It reads all key/value pairs from a source path and writes them to a
// destination path, optionally restricting the operation to a named subset
// of keys and supporting a dry-run mode that reports changes without
// persisting them.
package clone
