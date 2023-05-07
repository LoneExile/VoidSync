package minio

import (
	"voidsync/storage"
	"voidsync/sync"
)

type ginAPI struct {
	storageClient storage.Storage
	syncClient    sync.Syncer
}

func NewAPI(client storage.Storage, syncer sync.Syncer) *ginAPI {
	return &ginAPI{
		storageClient: client,
		syncClient:    syncer,
	}
}

func (api *ginAPI) GetRemoteFileList(remotePath string) (map[string]storage.FileInfo, error) {
	remoteFiles, err := api.storageClient.GetRemoteFileList(remotePath)
	if err != nil {
		return nil, err
	}
	return remoteFiles, nil
}

func (api *ginAPI) Sync(localPath string, remotePath string) error {
	err := api.syncClient.Sync(api.storageClient, localPath, remotePath)
	if err != nil {
		return err
	}
	return nil
}
