
# storage-sync


A backup and sync tool for local filesystems and Minio (S3-compatible) object storage.

**Note:** could be separate repo, but keeping this here for now, as the application is quite straightfoward.

## Features
- **Backup/Sync**: Copies files from a local source directory to a target (local directory or Minio bucket).
- **No Deletion**: Files deleted locally are NOT deleted from the target. The target acts as a backup/archive.
- **Modes**:
  - `distcp`: Full copy, overwrites all files in the target.
  - `sync`: Only copies new or changed files (by mod time or presence).
  - `sync` can be run once or continuously at a configurable interval.
- **Backends**:
  - Local filesystem
  - Minio (S3-compatible)
- **CLI**: Command-line interface for easy automation and scripting.

## Usage


### CLI Example
```
# Run directly with Go
go run ./cmd/cli.go --backend fs --src /path/to/source --trg /path/to/target --mode sync
go run ./cmd/cli.go --backend minio --src /path/to/source --trg backup/ --bucket mybucket --mode sync --interval 60


# Build the CLI binary
go build -o storage-sync-cli ./storage-sync/cmd/cli.go


# Or using make commands:
make build


# Run the built binary
./storage-sync-cli --backend fs --src /path/to/source --trg /path/to/target --mode sync
./storage-sync-cli --backend minio --src /path/to/source --trg backup/ --bucket mybucket --mode sync --interval 60
# Using make
make run --backend minio --src /path/to/source --trg backup/ --bucket mybucket --mode sync
```

- `--backend`: `fs` (local filesystem) or `minio`
- `--src`: Source directory (required)
- `--trg`: Target directory (fs) or prefix (minio)
- `--bucket`: Minio bucket name (required for minio backend)
- `--mode`: `sync` (default) or `distcp`
- `--interval`: If set, runs sync continuously every N seconds

### Environment Variables for Minio
- `MINIO_ENDPOINT`: Minio server endpoint (e.g., `localhost:9000`)
- `MINIO_ACCESS_KEY`: Minio access key
- `MINIO_SECRET_KEY`: Minio secret key
- `MINIO_USE_SSL`: `true` or `false`

## Design Notes
- This tool is designed for backup and one-way sync. It will **not** delete files from the target if they are deleted locally.
- If you want true mirroring (including deletions), you must implement or enable that feature explicitly.
- All operations are logged to stdout/stderr.

## Extending
- Add new backends by implementing the `WatcherBackend` interface.
- Add new CLI options or modes as needed.

## License
MIT
