package main

import (
	"context"
	"fmt"

	"github.com/CocaineCong/tangseng/app/index_platform/analyzer"
	"github.com/CocaineCong/tangseng/app/index_platform/cmd/kfk_register"
	"github.com/CocaineCong/tangseng/app/index_platform/repository/storage"
	"github.com/CocaineCong/tangseng/app/index_platform/rpc"
	"github.com/CocaineCong/tangseng/app/index_platform/trie"
	"github.com/CocaineCong/tangseng/loading"
)

func main() {
	ctx := context.Background()
	loading.Loading()
	analyzer.InitSeg()
	rpc.Init()
	kfk_register.RegisterJob(ctx)
	trie.InitTrieTree()
	storage.InitTrieDBs()

	fmt.Println("worker 开始工作...")
	fmt.Println("worker 结束工作...")
}
