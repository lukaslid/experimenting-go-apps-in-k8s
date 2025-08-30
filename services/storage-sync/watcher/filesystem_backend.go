package watcher

import (
	"os"
	"strings"
)

type FileSystemBackend struct {
	SrcPrefix string
	TrgPrefix string
}

func (b *FileSystemBackend) Initialize() error {
	// Initialize creates the target directory if it does not exist.
	// Returns an error if the directory cannot be created.
	return os.MkdirAll(b.TrgPrefix, 0755)
}

func (b *FileSystemBackend) DistCp() error {
	// DistCp copies all files and directories from SrcPrefix to TrgPrefix, overwriting existing files.
	// Returns an error if any file or directory cannot be copied.
	if err := os.MkdirAll(b.TrgPrefix, 0755); err != nil {
		return err
	}
	return WalkSourceFiles(b.SrcPrefix, func(path string) error {
		trgPath := strings.ReplaceAll(path, b.SrcPrefix, b.TrgPrefix)
		return b.CopyFile(path, trgPath)
	}, func(dir string) error {
		if dir != b.SrcPrefix {
			trgDir := strings.ReplaceAll(dir, b.SrcPrefix, b.TrgPrefix)
			return os.MkdirAll(trgDir, 0755)
		}
		return nil
	})
}

func (b *FileSystemBackend) CopyFile(srcPath string, trgPath string) error {
	// CopyFile copies a single file from srcPath to trgPath on the local filesystem.
	// Returns an error if the file cannot be copied.
	f, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer f.Close()

	w, err := os.Create(trgPath)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.ReadFrom(f)
	return err
}

// Sync copies only files that are new or have changed (by mod time) from source to target.
func (b *FileSystemBackend) Sync() error {
	// Sync copies only files that are new or have changed (by mod time) from source to target.
	// It skips files that are already up-to-date in the target directory.
	return WalkSourceFiles(b.SrcPrefix, func(path string) error {
		trgPath := strings.ReplaceAll(path, b.SrcPrefix, b.TrgPrefix)
		srcInfo, err := os.Stat(path)
		if err != nil {
			return err
		}
		trgInfo, err := os.Stat(trgPath)
		if err == nil {
			// If file exists and mod time is same or newer, skip
			if !srcInfo.ModTime().After(trgInfo.ModTime()) {
				return nil
			}
		}
		return b.CopyFile(path, trgPath)
	}, func(dir string) error {
		if dir != b.SrcPrefix {
			trgDir := strings.ReplaceAll(dir, b.SrcPrefix, b.TrgPrefix)
			return os.MkdirAll(trgDir, 0755)
		}
		return nil
	})
}
