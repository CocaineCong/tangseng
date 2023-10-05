package redis

import (
	"context"
	"fmt"
	"testing"

	"github.com/RoaringBitmap/roaring"

	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../config/config.yaml"}
	config.InitConfigForTest(&re)
	log.InitLog()
	InitRedis()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

func TestPushInvertedPath(t *testing.T) {
	ctx := context.Background()
	_ = PushInvertedPath(ctx, InvertedIndexDbPathDayKey, []string{"a", "b", "c"})
}

func TestListInvertedPath(t *testing.T) {
	ctx := context.Background()
	paths, _ := ListInvertedPath(ctx, []string{InvertedIndexDbPathDayKey})
	fmt.Println(paths)
}

func TestSetInvertedIndexTokenDocIds(t *testing.T) {
	ctx := context.Background()
	docIds := roaring.NewBitmap()
	docIds.AddInt(1)
	docIds.AddInt(2)
	err := SetInvertedIndexTokenDocIds(ctx, "test", docIds)
	fmt.Println(err)
}

func TestGetInvertedIndexTokenDocIds(t *testing.T) {
	ctx := context.Background()
	docIds, err := GetInvertedIndexTokenDocIds(ctx, "test1")
	fmt.Println(err)
	fmt.Println(docIds)
}

func TestPushInvertedIndexToken(t *testing.T) {
	ctx := context.Background()
	err := PushInvertedIndexToken(ctx, 1, "test2")
	fmt.Println(err)
}

func TestGetInvertedIndexToken(t *testing.T) {
	ctx := context.Background()
	tokens, err := ListInvertedIndexToken(ctx, 1)
	fmt.Println(err)
	fmt.Println(tokens)
}

func TestSetInvertedIndexByKey(t *testing.T) {
	ctx := context.Background()
	key := GetInvertedIndexDbPathMonthKey("10")
	SetInvertedPath(ctx, key, "a1")
	key2 := GetInvertedIndexDbPathMonthKey("11")
	SetInvertedPath(ctx, key2, "b")
}

func TestListInvertedIndexByPrefixKey(t *testing.T) {
	ctx := context.Background()
	key := GetInvertedIndexDbPathMonthKey("*")
	result, _ := ListInvertedIndexByPrefixKey(ctx, key)
	fmt.Println(result)
}
