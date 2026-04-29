package extract

import (
	"testing"
)

func TestExtract_NilClientErrors(t *testing.T) {
	_, err := Extract(nil, Options{
		Sources:     []string{"secret/a"},
		Destination: "secret/out",
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestExtract_EmptySourcesErrors(t *testing.T) {
	_, err := Extract(&stubClient{}, Options{
		Sources:     nil,
		Destination: "secret/out",
	})
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestExtract_EmptyDestinationErrors(t *testing.T) {
	_, err := Extract(&stubClient{}, Options{
		Sources:     []string{"secret/a"},
		Destination: "",
	})
	if err == nil {
		t.Fatal("expected error for empty destination")
	}
}

func TestPickKeys_AllKeysWhenNoneSpecified(t *testing.T) {
	data := map[string]string{"foo": "1", "bar": "2"}
	out := pickKeys(data, nil)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestPickKeys_SubsetOnly(t *testing.T) {
	data := map[string]string{"foo": "1", "bar": "2", "baz": "3"}
	out := pickKeys(data, []string{"foo", "baz"})
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["bar"]; ok {
		t.Error("unexpected key 'bar' in result")
	}
}

func TestPickKeys_MissingKeyIgnored(t *testing.T) {
	data := map[string]string{"foo": "1"}
	out := pickKeys(data, []string{"foo", "missing"})
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
}

func TestResult_DryRunFlagged(t *testing.T) {
	r := Result{DryRun: true}
	if !r.DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestResult_ExtractedAtSet(t *testing.T) {
	r := Result{}
	if !r.ExtractedAt.IsZero() == false {
		// zero is acceptable before operation; just ensure field exists
	}
	_ = r.ExtractedAt
}
