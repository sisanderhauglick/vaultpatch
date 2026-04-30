// Package broadcast provides functionality to fan-out a single Vault secret
// to one or more destination paths.
//
// Typical use-cases include propagating a shared credential (e.g. a database
// password) to every service path that needs it, or mirroring a config secret
// across multiple environment namespaces in a single operation.
//
// The caller may restrict which keys are propagated via Options.Keys; when the
// slice is empty every key present in the source secret is written.
package broadcast
