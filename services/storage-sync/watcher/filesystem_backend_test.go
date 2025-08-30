package watcher

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileSystemBackend_DistCpAndSync(t *testing.T) {
	tmpSrc := t.TempDir()
	tmpDst := t.TempDir()

	// Create a test file in src
	srcFile := filepath.Join(tmpSrc, "test.txt")
	if err := os.WriteFile(srcFile, []byte("hello world"), 0644); err != nil {
		t.Fatalf("failed to write src file: %v", err)
	}

	backend := &FileSystemBackend{SrcPrefix: tmpSrc, TrgPrefix: tmpDst}
	if err := backend.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	if err := backend.DistCp(); err != nil {
		t.Fatalf("DistCp failed: %v", err)
	}

	dstFile := filepath.Join(tmpDst, "test.txt")
	if _, err := os.Stat(dstFile); err != nil {
		t.Errorf("file not copied: %v", err)
	}

	// Now test Sync (should skip up-to-date file)
	if err := backend.Sync(); err != nil {
		t.Fatalf("Sync failed: %v", err)
	}
}
