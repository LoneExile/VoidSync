package minio

import (
	"time"
	"voidsync/storage"
	"voidsync/sync"
	"voidsync/utils"

	log "github.com/sirupsen/logrus"
)

func Sync(client storage.Storage, localPath, remotePath string) error {
	header := []string{"File", "Timestamp"}

	localFiles, err := sync.GetLocalFileList(localPath)
	if err != nil {
		log.Fatalln("Failed to get local file list:", err)
	}
	log.Info("\nLocal file list:")
	utils.LogTable(header, localFiles)

	remoteFiles, err := client.GetRemoteFileList(remotePath)
	if err != nil {
		log.Fatalln("Failed to get remote file list:", err)
	}
	log.Info("\nRemote file list:")
	utils.LogTable(header, remoteFiles)

	// ------------------------------------------
	remoteFileMap := make(map[string]time.Time)
	for filePath, timestamp := range remoteFiles {
		remoteFileMap[filePath] = timestamp
	}

	// Upload local files that are new or modified
	for filePath, localFileInfo := range localFiles {
		localPath := filePath
		localTimestamp := localFileInfo

		if remoteTimestamp, ok := remoteFileMap[localPath]; !ok || localTimestamp.After(remoteTimestamp) {
			log.Infof("Uploading %s", localPath)
			err = client.Upload(localPath)
			if err != nil {
				log.Errorf("Failed to upload %s: %v", localPath, err)
			}
		}
	}

	// Download remote files that are new or modified
	for filePath, remoteFileInfo := range remoteFiles {
		remotePath := filePath
		remoteTimestamp := remoteFileInfo

		if localTimestamp, ok := localFiles[remotePath]; !ok || remoteTimestamp.After(localTimestamp) {
			log.Infof("Downloading %s", remotePath)
			err = client.Download(remotePath, localPath)
			if err != nil {
				log.Errorf("Failed to download %s: %v", remotePath, err)
			}
		}
	}

	return nil
}
