package diff

import "github.com/your-org/vaultpatch/internal/vault"

// ChangeType describes the kind of change for a secret key.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Change represents a single key-level difference between two SecretMaps.
type Change struct {
	Key    string
	Type   ChangeType
	OldVal string
	NewVal string
}

// Diff computes the difference between a source and target SecretMap.
// Returns a slice of Change entries for keys that differ.
func Diff(src, dst vault.SecretMap) []Change {
	var changes []Change

	for k, srcVal := range src {
		dstVal, ok := dst[k]
		if !ok {
			changes = append(changes, Change{Key: k, Type: Removed, OldVal: srcVal})
		} else if srcVal != dstVal {
			changes = append(changes, Change{Key: k, Type: Modified, OldVal: srcVal, NewVal: dstVal})
		}
	}

	for k, dstVal := range dst {
		if _, ok := src[k]; !ok {
			changes = append(changes, Change{Key: k, Type: Added, NewVal: dstVal})
		}
	}

	return changes
}
