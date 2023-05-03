package minio

import (
	"context"

	"github.com/minio/minio-go/v7"
)

func (m *MinioStorage) StatObject(path string) (interface{}, error) {
	ctx := context.Background()
	bucketName := m.Cfg.MinIOBucketName
	object, err := m.Client.StatObject(ctx, bucketName, path, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	return object, nil
}
