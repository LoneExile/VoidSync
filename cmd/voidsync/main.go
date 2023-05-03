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
	localPath := "./public/upload/"
	remotePath := "/"

	err = sMinio.Sync(client, localPath, remotePath)
	if err != nil {
		fmt.Println("Error syncing:", err)
		return
	}
}
