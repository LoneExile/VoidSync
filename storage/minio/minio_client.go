package minio

import (
	"context"
	"fmt"
	"log"

	"voidsync/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func newMinioClient(cfg *config.Config) (*minio.Client, error) {
	minioClient, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKeyID, cfg.MinIOSecretKey, ""),
		Secure: cfg.MinIOUseSSL,
	})

	if err != nil {
		log.Println("🔴 Error initializing MinIO client object:", err)
		return nil, err
	}

	logMessage := fmt.Sprintf("✅ MinIO client object initialized: endpoint=%s, secure=%t", cfg.MinIOEndpoint, cfg.MinIOUseSSL)
	log.Println(logMessage)
	return minioClient, nil
}

func createBucket(minioClient *minio.Client, cfg *config.Config) error {
	ctx := context.Background()
	bucketName := cfg.MinIOBucketName
	log.Println("Creating bucket:", bucketName)

	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("🚧 We already own %s\n", bucketName)
		} else {
			return err
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	return nil
}
