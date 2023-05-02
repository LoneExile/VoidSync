package storage

import (
	"time"
	"voidsync/config"
)

type Storage interface {
	InitClient(cfg *config.Config) (Storage, error)
	CreateBucket() error
	Upload(serverPath string) error
	Download(serverPath string, localPath string) error
	GetRemoteTimestamp(path string) (time.Time, error)
	GetRemoteFileList(prefix string) (map[string]time.Time, error)
}
