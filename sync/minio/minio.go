package minio

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

	// Upload local files that are new or modified
	for filePath, localFileInfo := range localFiles {
		if remoteTimestamp, ok := remoteFiles[filePath]; !ok || localFileInfo.After(remoteTimestamp) {
			log.Infof("Uploading %s", filePath+"::"+filePath)
			err = client.UploadFile(localPath, filePath, "application/octet-stream")
			if err != nil {
				log.Errorf("Failed to upload %s: %v", localPath, err)
			}
		}
	}

	// Download remote files that are new or modified
	tmpDir := mkTmpDir()
	isDownload := false
	for filePath, remoteFileInfo := range remoteFiles {
		ctx := context.Background()
		downloadTmpPath := filepath.Join(tmpDir, filePath)
		downloadRemotePath := filepath.Join(remotePath, filePath)

		if localTimestamp, ok := localFiles[filePath]; !ok || remoteFileInfo.After(localTimestamp) {
			isDownload = true
			log.Infof("Downloading to %s", filepath.Join("/tmp/", filePath))
			err = client.DownloadObject(ctx, downloadRemotePath, downloadTmpPath)
			if err != nil {
				log.Errorf("Failed to download %s: %v", downloadRemotePath, err)
			}
		}
	}

	if isDownload {
		for filePath := range remoteFiles {
			downloadTempPath := filepath.Join(tmpDir, filePath)
			localFilePath := filepath.Join(localPath, filePath)
			println("Moving", filepath.Join("/tmp/", filePath), "to", localFilePath)

			destDir := filepath.Dir(localFilePath)
			if _, err := os.Stat(destDir); os.IsNotExist(err) {
				err := os.MkdirAll(destDir, os.ModePerm)
				if err != nil {
					log.Errorf("Failed to create directory %s: %v", destDir, err)
					continue
				}
				println("Created", destDir)
			}

			if _, err := os.Stat(downloadTempPath); err == nil {
				err = os.Rename(downloadTempPath, localFilePath)
				if err != nil {
					log.Errorf("Failed to move file : %v", err)
				}
			}
		}
	}

	return nil
}

func mkTmpDir() string {
	tempDir, err := ioutil.TempDir("", "voidsync")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
	}
	defer os.RemoveAll(tempDir)

	fmt.Println("Temporary directory created:", tempDir)

	return tempDir
}
