package sync

type Sync interface {
	StartSync() error
	IsSyncing() bool
}

type NowSync[T any] interface {
	Sync
	SyncNow() (T, error)
}
