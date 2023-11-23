package storage

import (
	"context"
	"fmt"
	"testing"
)

func TestInitInvertedDB(t *testing.T) {
	ctx := context.Background()
	InitInvertedDB(ctx)
	for _, v := range GlobalInvertedDB {
		fmt.Println(v)
	}
}
