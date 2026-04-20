// Package protect implements write-protection for Vault secret paths.
//
// A protected path has a metadata key injected into its secret map.
// Any vaultpatch command that writes secrets should check IsProtected
// before proceeding and abort with an informative error when true.
package protect
