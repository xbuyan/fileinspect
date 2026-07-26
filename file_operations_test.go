package main

import (
	"os"
	"path/filepath"
	"testing"
)

func withTempDir(t *testing.T, fn func(dir string)) {
	t.Helper()
	dir := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)
	fn(dir)
}

func TestHashFile_ProducesConsistentHash(t *testing.T) {
	withTempDir(t, func(dir string) {
		path := filepath.Join(dir, "sample.txt")
		if err := os.WriteFile(path, []byte("hello world"), 0o644); err != nil {
			t.Fatal(err)
		}

		h1, err := hashFile(path)
		if err != nil {
			t.Fatalf("hashFile failed: %v", err)
		}
		h2, err := hashFile(path)
		if err != nil {
			t.Fatalf("hashFile failed: %v", err)
		}
		if h1 != h2 {
			t.Fatalf("expected identical hashes for unchanged file, got %q and %q", h1, h2)
		}
		if h1 == "" {
			t.Fatal("expected non-empty hash")
		}
	})
}

func TestHashFile_DifferentContentDifferentHash(t *testing.T) {
	withTempDir(t, func(dir string) {
		pathA := filepath.Join(dir, "a.txt")
		pathB := filepath.Join(dir, "b.txt")
		os.WriteFile(pathA, []byte("hello"), 0o644)
		os.WriteFile(pathB, []byte("world"), 0o644)

		hA, err := hashFile(pathA)
		if err != nil {
			t.Fatal(err)
		}
		hB, err := hashFile(pathB)
		if err != nil {
			t.Fatal(err)
		}
		if hA == hB {
			t.Fatal("expected different content to produce different hashes")
		}
	})
}

func TestProcessSingleFile_NonExistentFile(t *testing.T) {
	withTempDir(t, func(dir string) {
		err := processSingleFile(filepath.Join(dir, "does-not-exist.txt"))
		if err == nil {
			t.Fatal("expected error for non-existent file, got nil")
		}
	})
}

func TestProcessSingleFile_RejectsDirectory(t *testing.T) {
	withTempDir(t, func(dir string) {
		sub := filepath.Join(dir, "subdir")
		os.Mkdir(sub, 0o755)

		err := processSingleFile(sub)
		if err == nil {
			t.Fatal("expected error when passing a directory to single-file mode, got nil")
		}
	})
}

func TestProcessSingleFile_SavesRecord(t *testing.T) {
	withTempDir(t, func(dir string) {
		path := filepath.Join(dir, "doc.txt")
		os.WriteFile(path, []byte("some content"), 0o644)

		if err := processSingleFile(path); err != nil {
			t.Fatalf("processSingleFile failed: %v", err)
		}

		if _, err := os.Stat("record.json"); err != nil {
			t.Fatalf("expected record.json to be created: %v", err)
		}
	})
}

func TestBatchHash_ProcessesAllFiles(t *testing.T) {
	withTempDir(t, func(dir string) {
		os.WriteFile(filepath.Join(dir, "a.txt"), []byte("aaa"), 0o644)
		os.WriteFile(filepath.Join(dir, "b.txt"), []byte("bbb"), 0o644)
		sub := filepath.Join(dir, "sub")
		os.Mkdir(sub, 0o755)
		os.WriteFile(filepath.Join(sub, "c.txt"), []byte("ccc"), 0o644)

		if err := batchHash(dir); err != nil {
			t.Fatalf("batchHash failed: %v", err)
		}

		data, err := os.ReadFile("batch_record.json")
		if err != nil {
			t.Fatalf("expected batch_record.json to be created: %v", err)
		}
		if len(data) == 0 {
			t.Fatal("expected batch_record.json to be non-empty")
		}
	})
}

func TestVerifyFile_UnmodifiedFilePasses(t *testing.T) {
	withTempDir(t, func(dir string) {
		path := filepath.Join(dir, "doc.txt")
		os.WriteFile(path, []byte("original content"), 0o644)

		if err := processSingleFile(path); err != nil {
			t.Fatalf("processSingleFile failed: %v", err)
		}

		if err := verifyFile(path); err != nil {
			t.Fatalf("expected verify to pass for unmodified file, got error: %v", err)
		}
	})
}

func TestVerifyFile_DetectsTamperedFile(t *testing.T) {
	withTempDir(t, func(dir string) {
		path := filepath.Join(dir, "doc.txt")
		os.WriteFile(path, []byte("original content"), 0o644)

		if err := processSingleFile(path); err != nil {
			t.Fatalf("processSingleFile failed: %v", err)
		}

		// Modify the file after recording — verify must now return an error
		// so the CLI exits non-zero and this is scriptable in CI/hooks.
		os.WriteFile(path, []byte("tampered content"), 0o644)

		if err := verifyFile(path); err == nil {
			t.Fatal("expected verify to return an error for a tampered file, got nil")
		}
	})
}
