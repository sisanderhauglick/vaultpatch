// Package patch provides targeted key-level mutations for Vault secret paths.
// Unlike a full write, patch merges or removes individual keys without
// overwriting unrelated fields in the same secret.
package patch
