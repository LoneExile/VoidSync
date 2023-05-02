package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"voidsync/config"
	"voidsync/storage"
	"voidsync/storage/minio"
	sMinio "voidsync/sync/minio"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)
}

func main() {
	cfg := config.LoadConfig()

	var store storage.Storage

	if cfg.StorageType == "minio" {
		store = minio.NewMinioStorage()
	} else {
		fmt.Println("Invalid storage type")
		return
	}

	client, err := store.InitClient(cfg)
	if err != nil {
		fmt.Println("Error initializing storage client:", err)
		return
	}

	err = client.CreateBucket()
	if err != nil {
		fmt.Println("Error creating bucket:", err)
		return
	}

	// ----------------------------------------------------------------------

	// // Upload a file or directory
	// path := "./public/upload"
	// err = client.Upload(path)
	// if err != nil {
	// 	log.Fatalln("Failed to upload:", err)
	// }

	// // Download a file
	// serverPath := "/"
	// localPath := "./public/download/"
	// err = client.Download(serverPath, localPath)
	// if err != nil {
	// 	log.Fatalln("Failed to download:", err)
	// }

	// ----------------------------------------------------------------------
	localPath := "./public/upload/"
	remotePath := "/"

	err = sMinio.Sync(client, localPath, remotePath)
	if err != nil {
		fmt.Println("Error syncing:", err)
		return
	}
}
