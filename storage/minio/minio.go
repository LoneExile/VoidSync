package minio

import (
	"context"
	"time"
	"voidsync/config"
	"voidsync/storage"
	"voidsync/utils"

	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	Client *minio.Client
	Cfg    *config.Config
}

func (m *MinioStorage) InitClient(cfg *config.Config) (storage.Storage, error) {
	minioClient, err := newMinioClient(cfg)
	if err != nil {
		return nil, err
	}

	m.Client = minioClient
	m.Cfg = cfg
	return m, nil
}

func (m *MinioStorage) CreateBucket() error {
	return createBucket(m.Client, m.Cfg)
}

func NewMinioStorage() storage.Storage {
	return &MinioStorage{}
}

func (m *MinioStorage) GetRemoteTimestamp(path string) (time.Time, error) {
	ctx := context.Background()
	bucketName := m.Cfg.MinIOBucketName
	object, err := m.Client.StatObject(ctx, bucketName, path, minio.StatObjectOptions{})
	if err != nil {
		return time.Time{}, err
	}

	return object.LastModified, nil
}

func (m *MinioStorage) GetRemoteFileList(prefix string) (map[string]time.Time, error) {
	ctx := context.Background()
	fileList := make(map[string]time.Time)

	objectCh := m.Client.ListObjects(ctx, m.Cfg.MinIOBucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		fileList[object.Key] = object.LastModified
	}

	// for object := range objectCh {
	// 	if object.Err != nil {
	// 		return nil, object.Err
	// 	}

	// 	objectName := object.Key

	// 	// Get the object's metadata
	// 	objectInfo, err := m.Client.StatObject(ctx, m.Cfg.MinIOBucketName, objectName, minio.StatObjectOptions{})
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	// Extract the original modification time from the metadata
	// 	originalModTimeStr := objectInfo.Metadata.Get("X-Amz-Meta-Original-Modtime")
	// 	originalModTime, err := time.Parse(time.RFC3339, originalModTimeStr)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	// Use the original modification time instead of object.LastModified
	// 	fileList[objectName] = originalModTime
	// }

	for filePath, timestamp := range fileList {
		fileList[filePath] = utils.ConvertTimestamp(timestamp)
	}

	return fileList, nil
}
