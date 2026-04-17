// Package promote implements secret promotion between Vault paths.
//
// Promotion reads secrets from a source path, computes a diff against the
// destination path, and optionally writes the merged result back. A dry-run
// mode allows previewing changes without mutating Vault.
//
// Example usage:
//
//	result, err := promote.Promote(client, promote.Options{
//		SrcPath: "secret/staging/app",
//		DstPath: "secret/production/app",
//		DryRun:  true,
//	})
package promote
