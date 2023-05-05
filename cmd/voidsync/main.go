package main

import (
	"log"
	"voidsync/config"
	"voidsync/db"
	"voidsync/storage"
	"voidsync/storage/minio"
	sMinio "voidsync/sync/minio"
)

func main() {

	cfg := config.LoadConfig()

	var store storage.Storage

	if cfg.StorageType == "minio" {
		store = minio.NewMinioStorage()
	} else {
		log.Println("🔴 Invalid storage type")
		return
	}

	client, err := store.InitClient(cfg)
	if err != nil {
		log.Println("🔴 Error initializing storage client:", err)
		return
	}

	err = client.CreateBucket()
	if err != nil {
		log.Println("🔴 Error creating bucket:", err)
		return
	}

	// ----------------------------------------------------------------------
	localPath := "./public/upload/"
	remotePath := "/"

	err = sMinio.Sync(client, localPath, remotePath)
	if err != nil {
		log.Println("🔴 Error syncing:", err)
		return
	}

	db.Init()
}
