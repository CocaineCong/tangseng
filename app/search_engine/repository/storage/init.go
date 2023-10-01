package storage

import (
	"context"
)

func InitStorageDB(ctx context.Context) {
	InitInvertedDB(ctx)
	InitGlobalTrieDB(ctx)
}
