package minio

import (
	"context"
	"os"
	"time"
	"voidsync/config"

	"github.com/minio/minio-go/v7"
)

func LocalFileModTime(filePath string) (time.Time, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}
	return fileInfo.ModTime(), nil
}

func ServerFileModTime(minioClient *minio.Client, cfg *config.Config, objectName string) (time.Time, error) {
	ctx := context.Background()
	bucketName := cfg.MinIOBucketName

	object, err := minioClient.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return time.Time{}, err
	}

	return object.LastModified, nil
}

