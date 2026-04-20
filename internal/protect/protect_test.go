package protect_test

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/protect"
)

func TestProtect_NilClientErrors(t *testing.T) {
	_, err := protect.Protect(nil, protect.Options{Paths: []string{"secret/a"}, Owner: "alice"})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestProtect_EmptyPathsError(t *testing.T) {
	_, err := protect.Protect(stubClient(), protect.Options{Owner: "alice"})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestProtect_EmptyOwnerErrors(t *testing.T) {
	_, err := protect.Protect(stubClient(), protect.Options{Paths: []string{"secret/a"}})
	if err == nil {
		t.Fatal("expected error for empty owner")
	}
}

func TestUnprotect_NilClientErrors(t *testing.T) {
	_, err := protect.Unprotect(nil, protect.Options{Paths: []string{"secret/a"}})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestUnprotect_EmptyPathsError(t *testing.T) {
	_, err := protect.Unprotect(stubClient(), protect.Options{})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestIsProtected_NilClientErrors(t *testing.T) {
	_, err := protect.IsProtected(nil, "secret/a")
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}
