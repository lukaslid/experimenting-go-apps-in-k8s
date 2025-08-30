package watcher

import (
	"io/fs"
	"log"
	"path/filepath"
	"time"
)

// WatcherBackend defines the methods required for a backend implementation.
type WatcherBackend interface {
	// Initialize prepares the backend for use (e.g., creates directories or connects to remote storage).
	Initialize() error
	// DistCp performs a full copy from source to target, overwriting all files.
	DistCp() error
	// Sync only copies/upload files that are new or have changed.
	Sync() error
	// CopyFile copies or uploads a single file from srcPath to trgPath.
	CopyFile(srcPath string, trgPath string) error
}

func WalkSourceFiles(srcPrefix string, onFile func(path string) error, onDir func(path string) error) error {
	// WalkSourceFiles walks the srcPrefix directory, calling onFile for each file and onDir for each directory.
	// onFile is called for every file, onDir for every directory. If either callback returns an error, walking stops and the error is returned.
	return filepath.WalkDir(srcPrefix, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error walking path %s: %v", path, err)
			return err
		}
		if d.IsDir() {
			return onDir(path)
		}
		return onFile(path)
	})
}

type Watcher struct {
	// Backend is the implementation used for file operations (local FS, Minio, etc).
	Backend WatcherBackend
}

// NewWatcher creates a new Watcher with the given backend (assumes backend is already initialized).
func NewWatcher(backend WatcherBackend) Watcher {
	return Watcher{
		Backend: backend,
	}
}

// InitBackendWithWatcher initializes the backend and returns a Watcher if successful, or an error otherwise.
func InitBackendWithWatcher(backend WatcherBackend) (*Watcher, error) {
	if err := backend.Initialize(); err != nil {
		return nil, err
	}
	w := &Watcher{Backend: backend}
	return w, nil
}

// Watch performs a full copy (DistCp) using the already-initialized backend.
func (w *Watcher) Watch() error {
	return w.Backend.DistCp()
}

// ContinuousSync runs Sync in a loop with the given interval. It returns only if Sync returns a fatal error.
func (w *Watcher) ContinuousSync(interval time.Duration) error {
	for {
		if err := w.Backend.Sync(); err != nil {
			return err
		}
		time.Sleep(interval)
	}
}
