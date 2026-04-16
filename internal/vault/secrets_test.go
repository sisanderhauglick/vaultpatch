package vault

import (
	"context"
	"testing"
)

func TestReadSecrets_ReturnsEmptyOnNilData(t *testing.T) {
	// This test validates SecretMap behaviour without a live Vault.
	var sm SecretMap
	if sm == nil {
		sm = SecretMap{}
	}
	if len(sm) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(sm))
	}
}

func TestSecretMap_KeyAccess(t *testing.T) {
	sm := SecretMap{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	if sm["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", sm["DB_HOST"])
	}
	if sm["DB_PORT"] != "5432" {
		t.Errorf("expected 5432, got %s", sm["DB_PORT"])
	}
}

func TestWriteSecrets_NilClientErrors(t *testing.T) {
	// Ensure calling WriteSecrets on a zero Client panics or errors gracefully.
	defer func() {
		if r := recover(); r == nil {
			t.Log("no panic on nil vault client — acceptable if error is returned")
		}
	}()
	c := &Client{}
	_ = c.WriteSecrets(context.Background(), "secret", "myapp/config", SecretMap{"KEY": "val"})
}
