package minio

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"voidsync/config"

	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
)

func Download(minioClient *minio.Client, bucketName, prefix string, targetDir string) error {
	ctx := context.Background()

	// List all objects with the given prefix
	objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	// Download each object and save it to the target directory
	for object := range objectCh {
		if object.Err != nil {
			return object.Err
		}

		// Download the object
		obj, err := minioClient.GetObject(ctx, bucketName, object.Key, minio.GetObjectOptions{})
		if err != nil {
			return err
		}

		// Create target file path
		targetPath := filepath.Join(targetDir, object.Key)

		// Create target file's parent directory if not exists
		err = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
		if err != nil {
			return err
		}

		// Create target file
		targetFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer targetFile.Close()

		// Write the downloaded object to the target file
		_, err = io.Copy(targetFile, obj)
		if err != nil {
			return err
		}
	}

	log.Infof("Successfully downloaded %s to %s", prefix, targetDir)

	return nil
}

func Upload(minioClient *minio.Client, cfg *config.Config, path string) error {
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
