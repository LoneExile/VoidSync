package sync

import "time"

type Synchronizer interface {
	GetLocalFileList(path string) (map[string]time.Time, error)
}
