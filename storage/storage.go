package storage

import (
	"context"
	"time"
	"voidsync/config"
)

type Storage interface {
	InitClient(cfg *config.Config) (Storage, error)
	CreateBucket() error

	UploadFile(baseDir, filePath, contentType string) error

	DownloadObject(ctx context.Context, objectKey, targetDir string) error

	GetRemoteTimestamp(path string) (time.Time, error)
	GetRemoteFileList(prefix string) (map[string]time.Time, error)

	// helper functions
	StatObject(path string) (interface{}, error)
}
