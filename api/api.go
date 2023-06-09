package api

import (
	"voidsync/api/minio"
	"voidsync/storage"
	"voidsync/sync"
)

type API interface {
	GetRemoteFileList(remotePath string) (map[string]storage.FileInfo, error)
	Sync(localPath string, remotePath string) error

	DownloadObjectsInServer(localPath, remotePath string) error
	DownloadAllObjects(remotePath string, removeIcon bool) (string, error)

	UploadDirClient(localPath, remotePath, contentType string) error
}

func NewAPI(client storage.Storage, syncer sync.Syncer) API {
	return minio.NewAPI(client, syncer)
}
