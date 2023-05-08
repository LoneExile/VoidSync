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

		targetPath := targetDir
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
func (m *MinioStorage) DownloadAllObjects(ctx context.Context, prefix string) (string, error) {
	bucketName := m.Cfg.MinIOBucketName
	maxDownloadAttempts := m.Cfg.MaxDownloadAttempts

	objectCh := m.Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	tmpDir := utils.MkTmpDir()

	numWorkers := 5
	tasks := make(chan string, numWorkers)
	results := make(chan error, numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			for objectKey := range tasks {
				tmpFile := filepath.Join(tmpDir, objectKey)

				var err error
				for attempt := 0; attempt < maxDownloadAttempts; attempt++ {
					err = m.DownloadObject(ctx, objectKey, tmpFile)
					if err == nil {
						break
					}
					time.Sleep(1 * time.Second)
				}
				results <- err
			}
		}()
	}

	go func() {
		for object := range objectCh {
			if object.Err != nil {
				results <- object.Err
			} else {
				tasks <- object.Key
			}
		}
		close(tasks)
	}()

	var errCount int
	for i := 0; i < len(objectCh); i++ {
		err := <-results
		if err != nil {
			log.Printf("🔴 Failed to download object: %v", err)
			errCount++
		}
	}

	if errCount > 0 {
		return "", fmt.Errorf("failed to download %d objects", errCount)
	}

	log.Println("✅ Successfully downloaded all objects")

	return tmpDir, nil
}

// TODO: Add a progress bar.
func (m *MinioStorage) DownloadObjectsInServer(ctx context.Context, prefix, targetDir string) error {

	tmpDir, err := m.DownloadAllObjects(ctx, prefix)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	if err := utils.MoveFiles(tmpDir, targetDir); err != nil {
		log.Println("🔴 failed to move files from tmp dir to target dir")
		return err
	}

	return nil
}
