package sync

import "voidsync/storage"

type Syncer interface {
	Sync(client storage.Storage, localPath, remotePath string) error
}

