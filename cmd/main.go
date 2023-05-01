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
	// err = minio.Upload(minioClient, cfg, path)
	// if err != nil {
	// 	log.Fatalln("Failed to upload:", err)
	// }

	// // Download a file
	// serverPath := "/"
	// localPath := "./public/download"
	// err = minio.DownloadObjects(minioClient, cfg, serverPath, localPath)
	// if err != nil {
	// 	log.Fatalln("Failed to download:", err)
	// }

	// compare local and server file mod time
	localModTime, err := minio.LocalFileModTime(path)
	if err != nil {
		log.Fatalln("Failed to get local file mod time:", err)
	}
	log.Printf("Local file mod time: %v\n", localModTime)

	objectName := "test.txt"
	serverModTime, err := minio.ServerFileModTime(minioClient, cfg, objectName)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"bucketName": cfg.MinIOBucketName,
			"objectName": objectName,
		}).Fatalln("Failed to get server file mod time")
	}

	// Convert serverModTime to local time zone
	localZone := localModTime.Location()
	serverModTimeLocal := serverModTime.In(localZone)

	log.WithFields(log.Fields{
		"\nserverModTime":                      serverModTimeLocal,
		"\nlocalModTime":                       localModTime,
		"\nserverModTime.Equal(localModTime)":  serverModTimeLocal.Equal(localModTime),
		"\nserverModTime.After(localModTime)":  serverModTimeLocal.After(localModTime),
		"\nserverModTime.Before(localModTime)": serverModTimeLocal.Before(localModTime),
		"\nobjectName":                         objectName,
	}).Info("\n\n -- Compare server and local file mod time -- ")

}
