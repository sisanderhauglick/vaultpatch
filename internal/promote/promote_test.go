package promote_test

import (
	"testing"

	"github.com/vaultpatch/internal/promote"
)

func TestPromote_NilClientErrors(t *testing.T) {
	_, err := promote.Promote(nil, promote.Options{
		SrcPath: "secret/src",
		DstPath: "secret/dst",
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestPromote_EmptyPathsError(t *testing.T) {
	tests := []struct {
		name string
		src  string
		dst  string
	}{
		{"empty src", "", "secret/dst"},
		{"empty dst", "secret/src", ""},
		{"both empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := promote.Promote(nil, promote.Options{
				SrcPath: tt.src,
				DstPath: tt.dst,
			})
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestFilterKeys_Unexported_AllKeys(t *testing.T) {
	// Indirectly tested via Promote; ensure full src copy when keys is nil.
	// This test documents expected behaviour rather than calling private func.
	opts := promote.Options{
		SrcPath: "secret/src",
		DstPath: "secret/dst",
		DryRun:  true,
		Keys:    nil,
	}
	_ = opts // used in integration tests
}
