package minio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
	"voidsync/utils"

	"github.com/minio/minio-go/v7"
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
		return errors.New("🔴 failed to download object after multiple attempts: " + objectKey)
	}

	logMessage := fmt.Sprintf("✅ Successfully downloaded object: objectKey=%s, targetDir=%s", objectKey, targetDir)
	log.Println(logMessage)

	return nil
}

// TODO: Add a progress bar, download multiple objects concurrently using goroutines.
func (m *MinioStorage) DownloadAllObjects(ctx context.Context, prefix, targetDir string) error {
	bucketName := m.Cfg.MinIOBucketName
	maxDownloadAttempts := m.Cfg.MaxDownloadAttempts

	objectCh := m.Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return err
	}
	tmpDir := utils.MkTmpDir()
	defer os.RemoveAll(tmpDir)

	for object := range objectCh {
		if object.Err != nil {
			return object.Err
		}

		var err error
		for i := 0; i < maxDownloadAttempts; i++ {
			tmpFile := filepath.Join(tmpDir, object.Key)
			err = m.DownloadObject(ctx, object.Key, tmpFile)
			if err == nil {
				break
			}

			time.Sleep(1 * time.Second)
		}

		if err != nil {
			return errors.New("🔴 failed to download object after multiple attempts: " + object.Key)
		}
	}

	if err := utils.MoveFiles(tmpDir, targetDir); err != nil {
		log.Println("🔴 failed to move files from tmp dir to target dir")
		return err
	}
	log.Println("✅ Successfully downloaded all objects")

	return nil
}
