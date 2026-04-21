package sanitize_test

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/sanitize"
	"github.com/your-org/vaultpatch/internal/vault"
)

func TestSanitize_NilClientErrors(t *testing.T) {
	_, err := sanitize.Sanitize(nil, []string{"secret/a"}, sanitize.Options{})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestSanitize_EmptyPathsError(t *testing.T) {
	client := &vault.Client{}
	_, err := sanitize.Sanitize(client, nil, sanitize.Options{})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestShouldProcess_AllKeysWhenNoneSpecified(t *testing.T) {
	// indirectly tested via DryRun: every key in the map should appear in Changed
	client, secrets := stubClient(map[string]interface{}{
		"KEY": "  hello  ",
	})
	res, err := sanitize.Sanitize(client, []string{"secret/test"}, sanitize.Options{
		TrimSpace: true,
		DryRun:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = secrets
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if v, ok := res[0].Changed["KEY"]; !ok || v != "hello" {
		t.Errorf("expected Changed[KEY]=hello, got %q", v)
	}
}

func TestSanitize_RemoveEmpty(t *testing.T) {
	client, _ := stubClient(map[string]interface{}{
		"present": "value",
		"empty":   "",
	})
	res, err := sanitize.Sanitize(client, []string{"secret/test"}, sanitize.Options{
		RemoveEmpty: true,
		DryRun:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res[0].Removed) != 1 || res[0].Removed[0] != "empty" {
		t.Errorf("expected 'empty' in Removed, got %v", res[0].Removed)
	}
}

func TestSanitize_LowercaseKeys(t *testing.T) {
	client, _ := stubClient(map[string]interface{}{
		"DB_HOST": "localhost",
	})
	res, err := sanitize.Sanitize(client, []string{"secret/test"}, sanitize.Options{
		LowercaseKeys: true,
		DryRun:        true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res[0].Changed["db_host"]; !ok {
		t.Errorf("expected Changed to contain 'db_host', got %v", res[0].Changed)
	}
}

func TestSanitize_DryRun_FlaggedOnResult(t *testing.T) {
	client, _ := stubClient(map[string]interface{}{"k": " v "})
	res, _ := sanitize.Sanitize(client, []string{"secret/test"}, sanitize.Options{
		TrimSpace: true,
		DryRun:    true,
	})
	if !res[0].DryRun {
		t.Error("expected DryRun flag to be true")
	}
}

func TestSanitize_SanitizedAtSet(t *testing.T) {
	client, _ := stubClient(map[string]interface{}{"k": "v"})
	res, _ := sanitize.Sanitize(client, []string{"secret/test"}, sanitize.Options{DryRun: true})
	if res[0].SanitizedAt.IsZero() {
		t.Error("expected SanitizedAt to be set")
	}
}
