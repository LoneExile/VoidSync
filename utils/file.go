package utils

import (
	"os"
	"path/filepath"
	"time"
	"voidsync/storage"
)

func GetLocalFileList(path string) (map[string]storage.FileInfo, error) {
	fileList := make(map[string]storage.FileInfo)

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(path, filePath)
			if err != nil {
				return err
			}
			checksum, err := calculateMD5(filePath)
			if err != nil {
				return err
			}
			fileList[relPath] = storage.FileInfo{
				Path:      relPath,
				Timestamp: ConvertTimestamp(info.ModTime()),
				ETag:      checksum,
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileList, nil
}

func GetLocalFileListTime(path string) (map[string]time.Time, error) {
	fileList := make(map[string]time.Time)

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(path, filePath)
			if err != nil {
				return err
			}
			fileList[relPath] = ConvertTimestamp(info.ModTime())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileList, nil
}
