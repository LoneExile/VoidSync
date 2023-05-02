package sync

import (
	"os"
	"path/filepath"
	"time"
)

func GetLocalTimestamp(path string) (time.Time, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return fileInfo.ModTime(), nil
}

func GetLocalFileList(path string) (map[string]time.Time, error) {
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
			fileList[relPath] = info.ModTime()
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileList, nil
}
