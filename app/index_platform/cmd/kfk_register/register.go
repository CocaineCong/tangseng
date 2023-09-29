package kfk_register

import (
	"context"
)

func RegisterJob(ctx context.Context) {
	newCtx := ctx
	go RunTireTree(newCtx)
	// newCtx = ctx
	// go RunInvertedIndex(newCtx)
}
