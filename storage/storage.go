package storage

import (
	"context"
	"io"
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

	UploadFileClient(ctx context.Context, filePath, contentType string, fileSize int64, file io.Reader) error
	UploadDirClient(ctx context.Context, baseDir, remotePath, contentType string) error

	DownloadObject(ctx context.Context, objectKey, targetDir string) error
	DownloadObjectsInServer(ctx context.Context, prefix, targetDir string) error
	DownloadAllObjects(ctx context.Context, prefix string) (string, error)

	GetRemoteTimestamp(path string) (time.Time, error)
	GetRemoteFileList(prefix string) (map[string]FileInfo, error)

	// helper functions
	StatObject(path string) (interface{}, error)
}
