package kfk_register

import (
	"context"
)

func RegisterJob(ctx context.Context) {
	newCtx := ctx
	// go RunTireTree(newCtx) // TODO:这个有点问题，后续优化再开启
	// newCtx = ctx
	go RunInvertedIndex(newCtx)
}
