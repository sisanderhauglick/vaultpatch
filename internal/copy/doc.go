// Package copy implements the secret copy feature for vaultpatch.
//
// It reads secrets from a source Vault path and writes them to a destination
// path, with optional key filtering and dry-run support.
//
// Usage:
//
//	res, err := copy.Copy(client, copy.Options{
//		SourcePath: "secret/data/staging/app",
//		DestPath:   "secret/data/prod/app",
//		Keys:       []string{"DB_HOST", "DB_PORT"},
//		DryRun:     true,
//	})
package copy
