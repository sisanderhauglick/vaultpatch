// Package split distributes keys from a single Vault secret path across
// multiple destination paths according to a caller-supplied assignment map.
//
// Typical usage:
//
//	results, err := split.Split(client, "secret/app", split.Options{
//		Assignments: map[string][]string{
//			"secret/app/db":  {"DB_HOST", "DB_PASS"},
//			"secret/app/api": {"API_KEY"},
//		},
//	})
package split
