package minio

import (
	"context"
	"os"
	// "path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
)

func (m *MinioStorage) UploadFile(baseDir, filePath, contentType string) error {
	err := m.checkIncompleteUploads("")
	if err != nil {
		log.Error("Failed to check and abort incomplete uploads:", err)
		return err
	}

	ctx := context.Background()
	bucketName := m.Cfg.MinIOBucketName

	// NOTE: RFC 3339 is an RFC standard for time strings,
	// which is to say how to represent a timestamp in textual form
	// e.g. 2020-01-01T00:00:00Z
	fileInfo, err := os.Stat(baseDir + filePath)
	if err != nil {
		log.Error("Failed to get file info:", err)
		return err
	}

	modTime := fileInfo.ModTime()
	opts := minio.PutObjectOptions{
		ContentType: contentType, UserMetadata: map[string]string{
			"x-amz-meta-original-modtime": modTime.Format(time.RFC3339),
		}}

	info, err := m.Client.FPutObject(ctx, bucketName, filePath, baseDir+filePath, opts)
	if err != nil {
		log.Error("Failed to upload file:", err)
		return err
	}

	log.WithFields(log.Fields{
		"objectName": filePath,
		"size":       info.Size,
	}).Info("Successfully uploaded file")

	return nil
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
