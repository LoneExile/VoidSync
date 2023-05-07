package minio

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"voidsync/storage"
	"voidsync/utils"
)

type MinioSyncer struct{}

func NewMinioSyncer() *MinioSyncer {
	return &MinioSyncer{}
}

func (m *MinioSyncer) Sync(client storage.Storage, localPath, remotePath string) error {
	localFiles, err := utils.GetLocalFileList(localPath)
	header := []string{"File", "ETag", "Timestamp"}

	// // TODO: don't use Fatalln, handle error properly
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
	downloadFiles(client, localFiles, remoteFiles, localPath, remotePath, tmpDir)
	uploadFiles(client, localFiles, remoteFiles, localPath)

	return nil
}

func uploadFiles(client storage.Storage, localFiles, remoteFiles map[string]storage.FileInfo, localPath string) {
	for filePath, localFileInfo := range localFiles {
		if remoteFileInfo, ok := remoteFiles[filePath]; !ok || (localFileInfo.ETag != remoteFileInfo.ETag && localFileInfo.Timestamp.After(remoteFileInfo.Timestamp)) {
			log.Printf("Uploading %s\n", filePath)
			err := client.UploadFile(localPath, filePath, "application/octet-stream")
			if err != nil {
				log.Printf("🔴 Failed to upload %s: %v", localPath, err)
			}
		}
	}
}

func downloadFiles(client storage.Storage, localFiles, remoteFiles map[string]storage.FileInfo, localPath, remotePath string, tmpDir string) {
	isDownload := false

	for filePath, remoteFileInfo := range remoteFiles {
		if localFileInfo, ok := localFiles[filePath]; !ok || (remoteFileInfo.ETag != localFileInfo.ETag && remoteFileInfo.Timestamp.After(localFileInfo.Timestamp)) {
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

func moveFiles(localFiles map[string]storage.FileInfo, localPath, tmpDir string) {
	for filePath := range localFiles {
		downloadTempPath := filepath.Join(tmpDir, filePath)
		localFilePath := filepath.Join(localPath, filePath)
		// log.Println("Moving", filepath.Join("/tmp/", filePath), "to", localFilePath)

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
