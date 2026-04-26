// Package stamp provides functionality to annotate HashiCorp Vault secret paths
// with arbitrary metadata key-value pairs. Annotations are written as reserved
// keys prefixed with "__" (e.g. __owner, __env, __stamped_at) so they can be
// identified and filtered independently from application secrets.
//
// Stamp supports dry-run mode, which computes and returns results without
// committing any writes to Vault.
package stamp
