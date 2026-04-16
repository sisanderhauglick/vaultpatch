package diff

// Type represents the kind of change between two secret maps.
type Type string

const (
	Added    Type = "added"
	Removed  Type = "removed"
	Modified Type = "modified"
)

// Entry describes a single key-level change.
type Entry struct {
	Key      string
	Type     Type
	OldValue string
	NewValue string
}

// Diff compares a source secret map against a target and returns the
// ordered list of changes needed to make source look like target.
func Diff(source, target map[string]string) []Entry {
	var entries []Entry

	for k, tv := range target {
		if sv, ok := source[k]; !ok {
			entries = append(entries, Entry{Key: k, Type: Added, NewValue: tv})
		} else if sv != tv {
			entries = append(entries, Entry{Key: k, Type: Modified, OldValue: sv, NewValue: tv})
		}
	}

	for k, sv := range source {
		if _, ok := target[k]; !ok {
			entries = append(entries, Entry{Key: k, Type: Removed, OldValue: sv})
		}
	}

	return entries
}
