package main

import (
	"log"
	"voidsync/config"
	"voidsync/db"
	"voidsync/server"
	"voidsync/storage"
	"voidsync/storage/minio"
	"voidsync/sync"
	sMinio "voidsync/sync/minio"
)

func main() {
	cfg := config.LoadConfig()

	var store storage.Storage
	var syncer sync.Syncer

	if cfg.StorageType == "minio" {
		store = minio.NewMinioStorage()
		syncer = sMinio.NewMinioSyncer()
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

	db.Init(cfg)

	server.StartServer(client, syncer, cfg)
}
