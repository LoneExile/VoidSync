package main

import (
	log "github.com/sirupsen/logrus"

	"voidsync/config"
	"voidsync/minio"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize MinIO client
	minioClient, err := minio.NewMinioClient(cfg)
	if err != nil {
		log.Fatalln("Failed to initialize MinIO client")
	}

	// Create a new bucket
	err = minio.CreateBucket(minioClient, cfg)
	if err != nil {
		log.Fatalln("Failed to create bucket:", err)
	}

	// Upload a file or directory
	path := "./public/upload"
	err = minio.Upload(minioClient, cfg, path)
	if err != nil {
		log.Fatalln("Failed to upload:", err)
	}

	// Download a file
	serverPath := "/test.txt"
	localPath := "./public/download"
	err = minio.Download(minioClient, cfg.MinIOBucketName, serverPath, localPath)
	if err != nil {
		log.Fatalln("Failed to download:", err)
	}

}
