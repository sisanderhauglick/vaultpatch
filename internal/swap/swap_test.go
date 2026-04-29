package swap_test

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/swap"
)

func TestSwap_NilClientErrors(t *testing.T) {
	_, err := swap.Swap(nil, "secret/a", "secret/b", swap.Options{})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestSwap_EmptySourceErrors(t *testing.T) {
	c := newStubClient()
	_, err := swap.Swap(c, "", "secret/b", swap.Options{})
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestSwap_EmptyDestErrors(t *testing.T) {
	c := newStubClient()
	_, err := swap.Swap(c, "secret/a", "", swap.Options{})
	if err == nil {
		t.Fatal("expected error for empty destination")
	}
}

func TestSwap_DryRun_ResultFlagged(t *testing.T) {
	c := newStubClient()
	c.data["secret/a"] = map[string]interface{}{"foo": "1"}
	c.data["secret/b"] = map[string]interface{}{"foo": "2"}

	results, err := swap.Swap(c, "secret/a", "secret/b", swap.Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].DryRun {
		t.Error("expected DryRun to be true")
	}
	// originals must be unchanged
	if c.data["secret/a"]["foo"] != "1" {
		t.Error("dry run must not modify source")
	}
	if c.data["secret/b"]["foo"] != "2" {
		t.Error("dry run must not modify destination")
	}
}

func TestSwap_LiveRun_ValuesExchanged(t *testing.T) {
	c := newStubClient()
	c.data["secret/a"] = map[string]interface{}{"key": "alpha"}
	c.data["secret/b"] = map[string]interface{}{"key": "beta"}

	_, err := swap.Swap(c, "secret/a", "secret/b", swap.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.data["secret/a"]["key"] != "beta" {
		t.Errorf("expected secret/a key=beta, got %v", c.data["secret/a"]["key"])
	}
	if c.data["secret/b"]["key"] != "alpha" {
		t.Errorf("expected secret/b key=alpha, got %v", c.data["secret/b"]["key"])
	}
}

func TestSwap_SubsetKeys_OnlySwapsSpecified(t *testing.T) {
	c := newStubClient()
	c.data["secret/a"] = map[string]interface{}{"x": "1", "y": "A"}
	c.data["secret/b"] = map[string]interface{}{"x": "2", "y": "B"}

	_, err := swap.Swap(c, "secret/a", "secret/b", swap.Options{Keys: []string{"x"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.data["secret/a"]["x"] != "2" {
		t.Errorf("expected x swapped in a")
	}
	if c.data["secret/a"]["y"] != "A" {
		t.Errorf("expected y unchanged in a")
	}
}

func TestResult_SwappedAtSet(t *testing.T) {
	c := newStubClient()
	c.data["secret/a"] = map[string]interface{}{"v": "1"}
	c.data["secret/b"] = map[string]interface{}{"v": "2"}

	results, _ := swap.Swap(c, "secret/a", "secret/b", swap.Options{DryRun: true})
	if results[0].SwappedAt.IsZero() {
		t.Error("expected SwappedAt to be set")
	}
}
