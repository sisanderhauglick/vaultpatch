package clone

import "testing"

func TestDryRun_ResultFlagged(t *testing.T) {
	res := Result{Source: "src", Destination: "dst", KeysCopied: 3, DryRun: true}
	if !res.DryRun {
		t.Fatal("expected DryRun to be true")
	}
	if res.KeysCopied != 3 {
		t.Fatalf("expected 3 keys copied, got %d", res.KeysCopied)
	}
}

func TestLiveRun_ResultFlagged(t *testing.T) {
	res := Result{Source: "src", Destination: "dst", KeysCopied: 2, DryRun: false}
	if res.DryRun {
		t.Fatal("expected DryRun to be false")
	}
}

func TestResult_SourceDestinationPreserved(t *testing.T) {
	res := Result{Source: "secret/dev", Destination: "secret/prod", KeysCopied: 5}
	if res.Source != "secret/dev" {
		t.Fatalf("unexpected source: %s", res.Source)
	}
	if res.Destination != "secret/prod" {
		t.Fatalf("unexpected destination: %s", res.Destination)
	}
}
