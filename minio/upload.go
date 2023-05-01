package minio

import (
	"context"
	"os"
	"path/filepath"
	"voidsync/config"

	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
)

func Upload(minioClient *minio.Client, cfg *config.Config, path string) error {
	err := checkIncompleteUploads(minioClient, cfg, "")
	if err != nil {
		log.Error("Failed to check and abort incomplete uploads:", err)
		return err
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Error("Failed to get file info:", err)
		return err
	}

	if fileInfo.IsDir() {
		return uploadDir(minioClient, cfg, path)
	} else {
		return uploadFile(minioClient, cfg, path, filepath.Base(path), "application/octet-stream")
	}
}

func uploadFile(minioClient *minio.Client, cfg *config.Config, filePath, objectName, contentType string) error {
	ctx := context.Background()
	bucketName := cfg.MinIOBucketName

	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
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

func uploadDir(minioClient *minio.Client, cfg *config.Config, dirPath string) error {
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

			err = uploadFile(minioClient, cfg, filePath, relPath, "application/octet-stream")
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func checkIncompleteUploads(minioClient *minio.Client, cfg *config.Config, objectPrefix string) error {
	ctx := context.Background()
	bucketName := cfg.MinIOBucketName

	listUploads := minioClient.ListIncompleteUploads(ctx, bucketName, objectPrefix, true)

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
		err := minioClient.RemoveIncompleteUpload(ctx, bucketName, upload.Key)
		if err != nil {
			log.Error("Failed to abort incomplete upload:", err)
			return err
		}
	}

	return nil
}
