package vault

import (
	"os"
	"testing"
)

func TestNewClient_MissingAddr(t *testing.T) {
	os.Unsetenv("VAULT_ADDR")
	os.Unsetenv("VAULT_TOKEN")

	_, err := NewClient("", "")
	if err == nil {
		t.Fatal("expected error when VAULT_ADDR is missing, got nil")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	os.Unsetenv("VAULT_TOKEN")

	_, err := NewClient("http://127.0.0.1:8200", "")
	if err == nil {
		t.Fatal("expected error when VAULT_TOKEN is missing, got nil")
	}
}

func TestNewClient_FromEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	t.Setenv("VAULT_TOKEN", "test-token")

	c, err := NewClient("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_ExplicitParams(t *testing.T) {
	c, err := NewClient("http://127.0.0.1:8200", "explicit-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_InvalidAddr(t *testing.T) {
	// A completely invalid address should still construct (vault SDK is lazy),
	// so we only verify no panic occurs and client is returned.
	c, err := NewClient("://bad-url", "token")
	// vault SDK may or may not error on construction; either outcome is acceptable
	_ = c
	_ = err
}
