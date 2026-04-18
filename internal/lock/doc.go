// Package lock implements advisory locking for Vault secret paths.
//
// A lock is stored as a secret at <path>/.lock containing the owner,
// acquisition time, and expiry. Lock and Unlock operations support
// dry-run mode, which simulates the operation without writing to Vault.
package lock
