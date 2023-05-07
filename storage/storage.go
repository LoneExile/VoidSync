package storage

import (
	"context"
	"time"
	"voidsync/config"
)

type FileInfo struct {
	Path      string
	Timestamp time.Time
	ETag      string
}

type Storage interface {
	InitClient(cfg *config.Config) (Storage, error)
	CreateBucket() error

	UploadFile(baseDir, filePath, contentType string) error

	DownloadObject(ctx context.Context, objectKey, targetDir string) error
	DownloadAllObjects(ctx context.Context, prefix, targetDir string) error

	GetRemoteTimestamp(path string) (time.Time, error)
	GetRemoteFileList(prefix string) (map[string]FileInfo, error)

	// helper functions
	StatObject(path string) (interface{}, error)
}
