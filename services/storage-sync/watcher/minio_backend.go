package watcher

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/minio/minio-go"
)

type MinioBackend struct {
	SrcPrefix string
	TrgBucket string
	TrgPrefix string
	Client    *minio.Client
}

func (m *MinioBackend) Initialize() error {
	// Initialize and creates a Minio client.
	// Returns an error if credentials are missing or the client cannot be created.

	endpoint := os.Getenv("MINIO_ENDPOINT")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	useSSLstr := os.Getenv("MINIO_USE_SSL")

	useSSL, err := strconv.ParseBool(useSSLstr)
	if err != nil {
		log.Fatalf("Invalid valio for MINIO_USE_SSL %s", useSSLstr)
		return err
	}
	// Check if essential variables are set
	if endpoint == "" || accessKeyID == "" || secretAccessKey == "" {
		log.Fatal("Required MinIO environment variables are not set.")
	}

	if endpoint == "" {
		return fmt.Errorf("MINIO_ENDPOINT environment variable is not set")
	}

	if accessKeyID == "" {
		return fmt.Errorf("MINIO_ACCESS_KEY environment variable is not set")
	}

	if secretAccessKey == "" {
		return fmt.Errorf("MINIO_SECRET_KEY environment variable is not set")
	}
	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Successfully connected to MinIO!")
	m.Client = minioClient

	return m.CreateBucket()
}

func (m *MinioBackend) CreateBucket() error {
	if m.Client == nil {
		return &minio.ErrorResponse{Message: "Minio client not initialized"}
	}
	err := m.Client.MakeBucket(m.TrgBucket, "")
	if err != nil {
		if minio.ToErrorResponse(err).Code == "BucketAlreadyOwnedByYou" {
			log.Printf("Bucket %s already exists", m.TrgBucket)
			return nil
		}
		return err
	}
	log.Printf("Bucket %s created successfully", m.TrgBucket)
	return nil
}

func (m *MinioBackend) DistCp() error {
	// DistCp uploads all files from SrcPrefix to the Minio bucket, overwriting existing objects.
	// Returns an error if any file cannot be uploaded.
	return WalkSourceFiles(m.SrcPrefix, func(path string) error {
		relPath := strings.TrimPrefix(path, m.SrcPrefix)
		trgPath := m.TrgPrefix + relPath
		return m.CopyFile(path, trgPath)
	}, func(dir string) error {
		// Optionally create folders in Minio if needed
		return nil
	})
}

func (m *MinioBackend) CopyFile(srcPath string, trgPath string) error {
	// CopyFile uploads a single file from srcPath to trgPath in the Minio bucket.
	// Returns an error if the file cannot be uploaded.
	if m.Client == nil {
		return &minio.ErrorResponse{Message: "Minio client not initialized"}
	}
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	// Remove leading slash if present in trgPath
	if len(trgPath) > 0 && trgPath[0] == '/' {
		trgPath = trgPath[1:]
	}

	// Upload file to Minio
	_, err = m.Client.PutObject(
		m.TrgBucket,
		trgPath,
		file,
		fileInfo.Size(),
		minio.PutObjectOptions{},
	)
	if err == nil {
		log.Printf("INFO: Uploaded %s to bucket %s as %s", srcPath, m.TrgBucket, trgPath)
	}
	return err
}

// Sync uploads only files that are new or have changed (by mod time) from source to Minio.
func (m *MinioBackend) Sync() error {
	// Sync uploads only files that are new or have changed from source to Minio.
	// It skips files that already exist in the Minio bucket (by key).
	existing := make(map[string]struct{})
	doneCh := make(chan struct{})
	defer close(doneCh)
	for obj := range m.Client.ListObjects(m.TrgBucket, m.TrgPrefix, true, doneCh) {
		if obj.Err != nil {
			continue
		}
		relPath := strings.TrimPrefix(obj.Key, m.TrgPrefix)
		existing[relPath] = struct{}{}
	}
	return WalkSourceFiles(m.SrcPrefix, func(path string) error {
		relPath := strings.TrimPrefix(path, m.SrcPrefix)
		trgPath := m.TrgPrefix + relPath
		if _, ok := existing[relPath]; ok {
			// TODO: Optionally compare hashes or mod times for more robust sync
			return nil
		}
		return m.CopyFile(path, trgPath)
	}, func(dir string) error {
		// Optionally create folders in Minio if needed
		return nil
	})
}
