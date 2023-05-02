package minio

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
)

func (m *MinioStorage) Upload(serverPath string) error {
	err := m.checkIncompleteUploads("")
	if err != nil {
		log.Error("Failed to check and abort incomplete uploads:", err)
		return err
	}

	fileInfo, err := os.Stat(serverPath)
	if err != nil {
		log.Error("Failed to get file info:", err)
		return err
	}

	if fileInfo.IsDir() {
		return m.uploadDir(serverPath)
	} else {
		return m.uploadFile(serverPath, filepath.Base(serverPath), "application/octet-stream")
	}
}

func (m *MinioStorage) uploadFile(filePath, objectName, contentType string) error {
	ctx := context.Background()
	bucketName := m.Cfg.MinIOBucketName

	// NOTE: RFC 3339 is an RFC standard for time strings,
	// which is to say how to represent a timestamp in textual form
	// e.g. 2020-01-01T00:00:00Z
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Error("Failed to get file info:", err)
		return err
	}
	modTime := fileInfo.ModTime()
	opts := minio.PutObjectOptions{
		ContentType: contentType, UserMetadata: map[string]string{
			"x-amz-meta-original-modtime": modTime.Format(time.RFC3339),
		}}

	info, err := m.Client.FPutObject(ctx, bucketName, objectName, filePath, opts)
	if err != nil {
		log.Error("Failed to upload file:", err)
		return err
	}

	log.WithFields(log.Fields{
		"objectName": objectName,
		"size":       info.Size,
	}).Info("Successfully uploaded file")

	return nil
}

func (m *MinioStorage) uploadDir(dirPath string) error {
	return filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Error("Error accessing path:", err)
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(dirPath, filePath)
			if err != nil {
				log.Error("Failed to get relative path:", err)
				return err
			}

			err = m.uploadFile(filePath, relPath, "application/octet-stream")
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (m *MinioStorage) checkIncompleteUploads(objectPrefix string) error {
	ctx := context.Background()
	bucketName := m.Cfg.MinIOBucketName

	listUploads := m.Client.ListIncompleteUploads(ctx, bucketName, objectPrefix, true)

	for upload := range listUploads {
		if upload.Err != nil {
			log.Error("Error listing incomplete uploads:", upload.Err)
			return upload.Err
		}

		log.WithFields(log.Fields{
			"objectName": upload.Key,
			"uploadID":   upload.UploadID,
		}).Warning("Incomplete upload found")

		// Abort the incomplete upload
		err := m.Client.RemoveIncompleteUpload(ctx, bucketName, upload.Key)
		if err != nil {
			log.Error("Failed to abort incomplete upload:", err)
			return err
		}
	}

	return nil
}
