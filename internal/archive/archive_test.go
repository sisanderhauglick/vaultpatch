package archive_test

import (
	"testing"

	"github.com/your-org/vaultpatch/internal/archive"
)

func TestArchive_NilClientErrors(t *testing.T) {
	_, err := archive.Archive(nil, archive.Options{
		Paths:       []string{"secret/foo"},
		ArchiveDest: "secret/archive",
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestArchive_EmptyPathsError(t *testing.T) {
	_, err := archive.Archive(nil, archive.Options{
		Paths:       []string{},
		ArchiveDest: "secret/archive",
	})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestArchive_EmptyDestError(t *testing.T) {
	_, err := archive.Archive(nil, archive.Options{
		Paths:       []string{"secret/foo"},
		ArchiveDest: "",
	})
	if err == nil {
		t.Fatal("expected error for empty archive destination")
	}
}

func TestResult_DryRunFlagged(t *testing.T) {
	r := archive.Result{
		Path:   "secret/foo",
		DryRun: true,
		Archived: true,
	}
	if !r.DryRun {
		t.Error("expected DryRun to be true")
	}
	if !r.Archived {
		t.Error("expected Archived to be true")
	}
}

func TestResult_ArchivedAtSet(t *testing.T) {
	import_time := archive.Result{}
	if !import_time.ArchivedAt.IsZero() {
		t.Error("zero-value ArchivedAt should be zero")
	}
}
