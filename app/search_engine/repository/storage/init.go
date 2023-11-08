package storage

import (
	"context"
	"fmt"
)

func InitStorageDB(ctx context.Context) {
	InitInvertedDB(ctx)
	fmt.Println("InitInvertedDB finish")
	InitGlobalTrieDB(ctx)
	fmt.Println("InitGlobalTrieDB finish")
}
