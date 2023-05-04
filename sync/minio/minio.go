package minio

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
	"voidsync/storage"
	"voidsync/sync"
	"voidsync/utils"
)

func Sync(client storage.Storage, localPath, remotePath string) error {
	localFiles, err := sync.GetLocalFileList(localPath)
	header := []string{"File", "Timestamp"}
	if err != nil {
		log.Fatalln("Failed to get local file list:", err)
	}
	utils.LogTable(header, localFiles)

	remoteFiles, err := client.GetRemoteFileList(remotePath)
	if err != nil {
		log.Fatalln("Failed to get remote file list:", err)
	}
	utils.LogTable(header, remoteFiles)

	tmpDir := mkTmpDir()
	uploadFiles(client, localFiles, remoteFiles, localPath)
	downloadFiles(client, localFiles, remoteFiles, localPath, remotePath, tmpDir)

	return nil
}

func uploadFiles(client storage.Storage, localFiles, remoteFiles map[string]time.Time, localPath string) {
	for filePath, localFileInfo := range localFiles {
		if remoteTimestamp, ok := remoteFiles[filePath]; !ok || localFileInfo.After(remoteTimestamp) {
			log.Printf("Uploading %s\n", filePath)
			err := client.UploadFile(localPath, filePath, "application/octet-stream")
			if err != nil {
				log.Printf("🔴 Failed to upload %s: %v", localPath, err)
			}
		}
	}
}

func downloadFiles(client storage.Storage, localFiles, remoteFiles map[string]time.Time, localPath, remotePath string, tmpDir string) {
	isDownload := false

	for filePath, remoteFileInfo := range remoteFiles {
		if localTimestamp, ok := localFiles[filePath]; !ok || remoteFileInfo.After(localTimestamp) {
			isDownload = true
			downloadTmpPath := filepath.Join(tmpDir, filePath)
			downloadRemotePath := filepath.Join(remotePath, filePath)

			log.Printf("Downloading to %s", filepath.Join("/tmp/", filePath))
			err := client.DownloadObject(context.Background(), downloadRemotePath, downloadTmpPath)
			if err != nil {
				log.Printf("🔴 Failed to download %s: %v", downloadRemotePath, err)
			}
		}
	}

	if isDownload {
		moveFiles(localFiles, localPath, tmpDir)
	}
}

func moveFiles(localFiles map[string]time.Time, localPath, tmpDir string) {
	for filePath := range localFiles {
		downloadTempPath := filepath.Join(tmpDir, filePath)
		localFilePath := filepath.Join(localPath, filePath)
		log.Println("Moving", filepath.Join("/tmp/", filePath), "to", localFilePath)

		destDir := filepath.Dir(localFilePath)
		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			err := os.MkdirAll(destDir, os.ModePerm)
			if err != nil {
				log.Printf("🔴 Failed to create directory %s: %v", destDir, err)
				continue
			}
			log.Println("Created", destDir)
		}

		if _, err := os.Stat(downloadTempPath); err == nil {
			err = os.Rename(downloadTempPath, localFilePath)
			if err != nil {
				log.Printf("🔴 Failed to move file : %v", err)
			}
		}
	}
}

func mkTmpDir() string {
	tempDir, err := ioutil.TempDir("", "voidsync")
	if err != nil {
		log.Println("Error creating temporary directory:", err)
	}
	defer os.RemoveAll(tempDir)

	log.Println("Temporary directory created:", tempDir)

	return tempDir
}
