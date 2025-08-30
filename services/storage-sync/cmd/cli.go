package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/lukaslid/storage-sync/watcher"
)

func main() {
	backendType := flag.String("backend", "fs", "Backend type: fs or minio")
	src := flag.String("src", "", "Source directory (required)")
	trg := flag.String("trg", "", "Target directory (fs) or prefix (minio)")
	bucket := flag.String("bucket", "", "Minio bucket name (for minio backend)")
	mode := flag.String("mode", "sync", "Mode: sync or distcp")
	interval := flag.Int("interval", 0, "Continuous sync interval in seconds (0 for one-shot)")
	flag.Parse()

	if *src == "" {
		fmt.Println("--src is required")
		os.Exit(1)
	}

	var backend watcher.WatcherBackend
	if *backendType == "minio" {
		if *bucket == "" {
			fmt.Println("--bucket is required for minio backend")
			os.Exit(1)
		}
		backend = &watcher.MinioBackend{
			SrcPrefix: *src,
			TrgPrefix: *trg,
			TrgBucket: *bucket,
		}
	} else {
		if *trg == "" {
			fmt.Println("--trg is required for fs backend")
			os.Exit(1)
		}
		backend = &watcher.FileSystemBackend{
			SrcPrefix: *src,
			TrgPrefix: *trg,
		}
	}

	w, err := watcher.InitBackendWithWatcher(backend)
	if err != nil {
		fmt.Printf("Failed to initialize backend: %v\n", err)
		os.Exit(1)
	}

	switch *mode {
	case "distcp":
		err = w.Watch()
	case "sync":
		if *interval > 0 {
			err = w.ContinuousSync(time.Duration(*interval) * time.Second)
		} else {
			err = backend.Sync()
		}
	default:
		fmt.Println("Unknown mode. Use sync or distcp.")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
