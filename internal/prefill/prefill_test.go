package prefill

import (
	"testing"
)

func TestPrefill_NilClientErrors(t *testing.T) {
	_, err := Prefill(nil, Options{
		Paths:    []string{"secret/app"},
		Defaults: map[string]string{"key": "val"},
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestPrefill_EmptyPathsError(t *testing.T) {
	c := stubClient()
	_, err := Prefill(c, Options{
		Paths:    nil,
		Defaults: map[string]string{"key": "val"},
	})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestPrefill_EmptyDefaultsError(t *testing.T) {
	c := stubClient()
	_, err := Prefill(c, Options{
		Paths:    []string{"secret/app"},
		Defaults: nil,
	})
	if err == nil {
		t.Fatal("expected error for empty defaults")
	}
}

func TestPrefill_DryRun_ResultFlagged(t *testing.T) {
	c := stubClient()
	res, err := Prefill(c, Options{
		Paths:    []string{"secret/app"},
		Defaults: map[string]string{"missing_key": "default_val"},
		DryRun:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if !res[0].DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestPrefill_DryRun_FilledKeysPopulated(t *testing.T) {
	c := stubClient()
	res, err := Prefill(c, Options{
		Paths:    []string{"secret/app"},
		Defaults: map[string]string{"alpha": "1", "beta": "2"},
		DryRun:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res[0].Filled) != 2 {
		t.Errorf("expected 2 filled keys, got %d", len(res[0].Filled))
	}
}

func TestPrefill_ExistingKeyNotOverwritten(t *testing.T) {
	c := stubClientWithData(map[string]interface{}{"existing": "original"})
	res, err := Prefill(c, Options{
		Paths:    []string{"secret/app"},
		Defaults: map[string]string{"existing": "new_val", "fresh": "added"},
		DryRun:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res[0].Filled["existing"]; ok {
		t.Error("existing key should not appear in Filled")
	}
	if _, ok := res[0].Filled["fresh"]; !ok {
		t.Error("fresh key should appear in Filled")
	}
}

func TestResult_PrefillledAtSet(t *testing.T) {
	c := stubClient()
	res, err := Prefill(c, Options{
		Paths:    []string{"secret/app"},
		Defaults: map[string]string{"k": "v"},
		DryRun:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res[0].PrefillledAt.IsZero() {
		t.Error("expected PrefillledAt to be set")
	}
}
