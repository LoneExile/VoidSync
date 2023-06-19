package utils

import (
	"archive/zip"
	"io"
	"log"
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

func MkTmpDir() string {
	tempDir := filepath.Join(os.TempDir(), "voidsync")

	log.Println("✅ Temporary directory created:", tempDir)

	return tempDir
}

func MoveFiles(srcDir, targetDir string) error {
	log.Println("Moving files from", srcDir, "to", targetDir)
	localFiles, err := GetLocalFileListTime(srcDir)
	if err != nil {
		return err
	}

	for key := range localFiles {
		srcPath := filepath.Join(srcDir, key)
		dstPath := filepath.Join(targetDir, key)

		destDir := filepath.Dir(dstPath)
		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			log.Println("Creating directory:", destDir)
			err := os.MkdirAll(destDir, os.ModePerm)
			if err != nil {
				return err
			}
		}

		err := os.Rename(srcPath, dstPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateZipArchive(writer io.Writer, srcDir string) error {
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(srcDir, path)
			if err != nil {
				return err
			}
			zipFile, err := zipWriter.Create(relPath)
			if err != nil {
				return err
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(zipFile, file)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
