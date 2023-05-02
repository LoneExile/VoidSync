package minio

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
)

// NOTE:
// Keep in mind that this approach assumes that the files in the target directory are not being accessed or modified by other processes during the download process.
// I need to ensure that the target directory remains consistent even when accessed by other processes, I might need to implement file locking or other concurrency control mechanisms.

// TODO: use ioutil.TempDir() instead of filepath.Join(localPath, ".tmp") to create a temporary directory
func (m *MinioStorage) Download(serverPath string, localPath string) error {
	ctx := context.Background()
	bucketName := m.Cfg.MinIOBucketName

	objectCh := m.Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    serverPath,
		Recursive: true,
	})

	// Create a temporary directory for downloads
	tempDir := filepath.Join(localPath, ".tmp")
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return err
	}
	log.Infof("Successfully created temp directory %s ...", tempDir)
	defer os.RemoveAll(tempDir)

	var objectKeys []string
	for object := range objectCh {
		if object.Err != nil {
			return object.Err
		}

		err := m.downloadObject(ctx, object.Key, tempDir)
		if err != nil {
			return err
		}
		objectKeys = append(objectKeys, object.Key)
	}
	log.Infof("Successfully downloaded %d objects", len(objectKeys))

	// Move downloaded files from the temporary directory to the target directory
	for _, objectKey := range objectKeys {
		tempPath := filepath.Join(tempDir, objectKey)
		targetPath := filepath.Join(localPath, objectKey)

		err = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
		if err != nil {
			return err
		}

		err = os.Rename(tempPath, targetPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MinioStorage) downloadObject(ctx context.Context, objectKey, targetDir string) error {
	bucketName := m.Cfg.MinIOBucketName
	maxDownloadAttempts := m.Cfg.MaxDownloadAttempts

	var err error
	for i := 0; i < maxDownloadAttempts; i++ {
		obj, err := m.Client.GetObject(ctx, bucketName, objectKey, minio.GetObjectOptions{})
		if err != nil {
			continue
		}

		targetPath := filepath.Join(targetDir, objectKey)
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

		// Sleep for a moment before retrying
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return errors.New("failed to download object after multiple attempts: " + objectKey)
	}
	log.WithFields(log.Fields{
		"objectKey": objectKey,
		"targetDir": targetDir,
	}).Info("Successfully downloaded object")

	return nil
}
