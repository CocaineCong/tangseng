package recall

import (
	"context"
	"fmt"
	"testing"

	"github.com/CocaineCong/tangseng/app/gateway/rpc"
	"github.com/CocaineCong/tangseng/app/search_engine/repository/storage"
	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/repository/redis"
)

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../../../config/config.yaml"}
	config.InitConfigForTest(&re)
	log.InitLog()
	redis.InitRedis()
	rpc.Init()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

func TestGetTrieTreeFromRedis(t *testing.T) {
	ctx := context.Background()
	storage.InitGlobalTrieDB(ctx)
	for _, v := range storage.GlobalTrieDB {
		tree, err := v.GetTrieTreeDict()
		if err != nil {
			fmt.Println("tree ", err)
		}
		tree.TraverseForRecall()
	}

}

func TestRecall_SearchVector(t *testing.T) {
	ctx := context.Background()
	r := NewRecall()
	queries := []string{"小岛秀夫", "导演"}
	res, err := r.SearchVector(ctx, queries)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
