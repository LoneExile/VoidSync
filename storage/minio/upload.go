package minio

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
)

func (m *MinioStorage) UploadFileClient(ctx context.Context, filePath, contentType string, fileSize int64, file io.Reader) error {
	bucketName := m.Cfg.MinIOBucketName

	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	info, err := m.Client.PutObject(ctx, bucketName, filePath, file, fileSize, opts)
	if err != nil {
		log.Printf("🔴 Failed to upload file: %v", err)
		return err
	}

	log.Printf("✅ Successfully uploaded file: %s, size: %d", filePath, info.Size)
	return nil
}

func (m *MinioStorage) UploadDirClient(ctx context.Context, baseDir, remotePath, contentType string) error {
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(baseDir, path)
			if err != nil {
				return err
			}

			remoteFilePath := filepath.Join(remotePath, relPath)

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			fileSize := info.Size()

			err = m.UploadFileClient(ctx, remoteFilePath, contentType, fileSize, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("🔴 Failed to upload all files in directory: %v", err)
		return err
	}

	log.Printf("✅ Successfully uploaded all files in directory: %s", baseDir)
	return nil
}

func (m *MinioStorage) UploadFile(baseDir, filePath, contentType string) error {
	err := m.checkIncompleteUploads("")
	if err != nil {
		log.Println("🔴 Failed to check and abort incomplete uploads:", err)
		return err
	}

	ctx := context.Background()
	bucketName := m.Cfg.MinIOBucketName

	// NOTE: RFC 3339 is an RFC standard for time strings,
	// which is to say how to represent a timestamp in textual form
	// e.g. 2020-01-01T00:00:00Z
	fileInfo, err := os.Stat(baseDir + filePath)
	if err != nil {
		log.Println("🔴 Failed to get file info:", err)
		return err
	}

	modTime := fileInfo.ModTime()
	opts := minio.PutObjectOptions{
		ContentType: contentType, UserMetadata: map[string]string{
			"x-amz-meta-original-modtime": modTime.Format(time.RFC3339),
		}}

	info, err := m.Client.FPutObject(ctx, bucketName, filePath, baseDir+filePath, opts)
	if err != nil {
		log.Println("🔴 Failed to upload file:", err)
		return err
	}

	logMessage := fmt.Sprintf("✅ Successfully uploaded file: %s, size: %d", filePath, info.Size)
	log.Println(logMessage)

	return nil
}

func (m *MinioStorage) checkIncompleteUploads(objectPrefix string) error {
	ctx := context.Background()
	bucketName := m.Cfg.MinIOBucketName

	listUploads := m.Client.ListIncompleteUploads(ctx, bucketName, objectPrefix, true)

	for upload := range listUploads {
		if upload.Err != nil {
			log.Println("🔴 Error listing incomplete uploads:", upload.Err)
			return upload.Err
		}

		logMessage := fmt.Sprintf("Incomplete upload found: %s, uploadID: %s", upload.Key, upload.UploadID)
		log.Println(logMessage)

		// Abort the incomplete upload
		err := m.Client.RemoveIncompleteUpload(ctx, bucketName, upload.Key)
		if err != nil {
			log.Println("🔴 Failed to abort incomplete upload:", err)
			return err
		}
	}

	return nil
}
