package normalize

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/vault"
)

func TestNormalize_NilClientErrors(t *testing.T) {
	_, err := Normalize(nil, Options{Paths: []string{"secret/a"}})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestNormalize_EmptyPathsError(t *testing.T) {
	client := &vault.Client{}
	_, err := Normalize(client, Options{})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestShouldProcess_AllKeysWhenNoneSpecified(t *testing.T) {
	if !shouldProcess("ANY_KEY", nil) {
		t.Error("expected true when keys list is nil")
	}
	if !shouldProcess("ANY_KEY", []string{}) {
		t.Error("expected true when keys list is empty")
	}
}

func TestShouldProcess_SubsetKeys(t *testing.T) {
	keys := []string{"HOST", "PORT"}
	if !shouldProcess("HOST", keys) {
		t.Error("expected HOST to be processed")
	}
	if shouldProcess("PASSWORD", keys) {
		t.Error("expected PASSWORD to be skipped")
	}
}

func TestApplyRules_TrimSpace(t *testing.T) {
	secrets := map[string]interface{}{
		"key": "  hello  ",
	}
	changes := applyRules(secrets, Options{TrimSpace: true})
	if changes["key"] != "hello" {
		t.Errorf("expected 'hello', got %q", changes["key"])
	}
}

func TestApplyRules_UppercaseValues(t *testing.T) {
	secrets := map[string]interface{}{
		"env": "production",
	}
	changes := applyRules(secrets, Options{UppercaseValues: true})
	if changes["env"] != "PRODUCTION" {
		t.Errorf("expected 'PRODUCTION', got %q", changes["env"])
	}
}

func TestApplyRules_LowercaseKeys(t *testing.T) {
	secrets := map[string]interface{}{
		"DB_HOST": "localhost",
	}
	changes := applyRules(secrets, Options{LowercaseKeys: true})
	if _, ok := changes["db_host"]; !ok {
		t.Errorf("expected lowercased key 'db_host' in changes, got %v", changes)
	}
}

func TestApplyRules_NoChangeSkipped(t *testing.T) {
	secrets := map[string]interface{}{
		"key": "value",
	}
	changes := applyRules(secrets, Options{TrimSpace: true})
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %v", changes)
	}
}
