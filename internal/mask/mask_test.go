package mask_test

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/mask"
)

func TestMask_AllKeys_WhenNoneSpecified(t *testing.T) {
	secrets := map[string]string{"foo": "bar", "baz": "qux"}
	r := mask.Mask("secret/app", secrets, mask.Options{})
	for k, v := range r.Masked {
		if v != "***" {
			t.Errorf("key %s: expected *** got %s", k, v)
		}
	}
	if r.Redacted != 2 {
		t.Errorf("expected 2 redacted, got %d", r.Redacted)
	}
}

func TestMask_SubsetKeys(t *testing.T) {
	secrets := map[string]string{"password": "s3cr3t", "username": "admin"}
	r := mask.Mask("secret/app", secrets, mask.Options{Keys: []string{"password"}})
	if r.Masked["password"] != "***" {
		t.Errorf("expected password masked")
	}
	if r.Masked["username"] != "admin" {
		t.Errorf("expected username unmasked")
	}
	if r.Redacted != 1 {
		t.Errorf("expected 1 redacted, got %d", r.Redacted)
	}
}

func TestMask_CustomPlaceholder(t *testing.T) {
	secrets := map[string]string{"token": "abc123"}
	r := mask.Mask("secret/app", secrets, mask.Options{Keys: []string{"token"}, Placeholder: "<REDACTED>"})
	if r.Masked["token"] != "<REDACTED>" {
		t.Errorf("expected custom placeholder, got %s", r.Masked["token"])
	}
}

func TestMask_CaseInsensitiveKey(t *testing.T) {
	secrets := map[string]string{"API_KEY": "xyz"}
	r := mask.Mask("secret/app", secrets, mask.Options{Keys: []string{"api_key"}})
	if r.Masked["API_KEY"] != "***" {
		t.Errorf("expected case-insensitive match to mask API_KEY")
	}
}

func TestMask_EmptySecrets(t *testing.T) {
	r := mask.Mask("secret/empty", map[string]string{}, mask.Options{})
	if r.Redacted != 0 {
		t.Errorf("expected 0 redacted, got %d", r.Redacted)
	}
}

func TestMask_PathPreserved(t *testing.T) {
	r := mask.Mask("secret/myapp", map[string]string{"k": "v"}, mask.Options{})
	if r.Path != "secret/myapp" {
		t.Errorf("expected path preserved, got %s", r.Path)
	}
}
