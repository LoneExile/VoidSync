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

func (m *MinioStorage) GetRemoteFileList(prefix string) (map[string]storage.FileInfo, error) {
	ctx := context.Background()
	fileList := make(map[string]storage.FileInfo)

	objectCh := m.Client.ListObjects(ctx, m.Cfg.MinIOBucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		fileList[object.Key] = storage.FileInfo{
			Path:      object.Key,
			Timestamp: utils.ConvertTimestamp(object.LastModified),
			ETag:      object.ETag,
		}
	}

	return fileList, nil
}
