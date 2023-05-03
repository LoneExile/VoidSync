package minio

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	// log "github.com/sirupsen/logrus"
)

// NOTE:
// Keep in mind that this approach assumes that the files in the target directory are not being accessed or modified by other processes during the download process.
// I need to ensure that the target directory remains consistent even when accessed by other processes, I might need to implement file locking or other concurrency control mechanisms.

func (m *MinioStorage) DownloadObject(ctx context.Context, objectKey, targetDir string) error {
	bucketName := m.Cfg.MinIOBucketName
	maxDownloadAttempts := m.Cfg.MaxDownloadAttempts

	var err error
	for i := 0; i < maxDownloadAttempts; i++ {

		obj, err := m.Client.GetObject(ctx, bucketName, objectKey, minio.GetObjectOptions{})
		if err != nil {
			continue
		}

		targetPath := filepath.Join(targetDir)
		err = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
		if err != nil {
			continue
		}

		targetFile, err := os.Create(targetPath)
		if err != nil {
			continue
		}
		defer targetFile.Close()

		_, err = io.Copy(targetFile, obj)
		if err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return errors.New("failed to download object after multiple attempts: " + objectKey)
	}
	// log.WithFields(log.Fields{
	// 	"objectKey": objectKey,
	// 	"targetDir": targetDir,
	// }).Info("Successfully downloaded object")

	return nil
}
