// Package patch provides targeted key-level mutations for Vault secret paths.
// Unlike a full write, patch merges or removes individual keys without
// overwriting unrelated fields in the same secret.
//
// Basic usage:
//
//	// Merge a single key into an existing secret
//	err := patch.Merge(client, "secret/data/myapp", map[string]interface{}{
//		"api_key": "new-value",
//	})
//
//	// Remove a key from an existing secret
//	err = patch.Remove(client, "secret/data/myapp", "api_key")
//
// Both operations read the current secret, apply the mutation, and write
// the result back, preserving all unmodified keys.
package patch
