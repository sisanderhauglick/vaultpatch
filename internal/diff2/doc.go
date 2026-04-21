// Package diff2 implements structured two-way secret diffing between
// two Vault paths. Unlike the basic diff package, diff2 returns typed
// Change records (added, removed, modified) with old and new values,
// supports key filtering, and records metadata such as the diff timestamp.
package diff2
