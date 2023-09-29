package service

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/RoaringBitmap/roaring"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/spf13/cast"

	"github.com/CocaineCong/tangseng/app/index_platform/analyzer"
	"github.com/CocaineCong/tangseng/app/index_platform/input_data"
	"github.com/CocaineCong/tangseng/consts/e"
	pb "github.com/CocaineCong/tangseng/idl/pb/index_platform"
	"github.com/CocaineCong/tangseng/pkg/mapreduce"
	"github.com/CocaineCong/tangseng/types"
)

type IndexPlatformSrv struct {
	*pb.UnimplementedIndexPlatformServiceServer
}

var (
	IndexPlatIns  *IndexPlatformSrv
	IndexPlatOnce sync.Once
)

func GetIndexPlatformSrv() *IndexPlatformSrv {
	IndexPlatOnce.Do(func() {
		IndexPlatIns = new(IndexPlatformSrv)
	})
	return IndexPlatIns
}

// BuildIndexService 构建索引
func (s *IndexPlatformSrv) BuildIndexService(ctx context.Context, req *pb.BuildIndexReq) (resp *pb.BuildIndexResp, err error) {
	// 时间估计
	resp = new(pb.BuildIndexResp)
	resp.Code = e.SUCCESS
	resp.Message = e.GetMsg(e.SUCCESS)
	invertedIndex := cmap.New[*roaring.Bitmap]() // 倒排索引

	// mapreduce 这个是用chan和goroutine来代替master和worker的rpc调用，避免了频繁的rpc调用
	_, _ = mapreduce.MapReduce(func(source chan<- []byte) {
		for _, path := range req.FilePath {
			content, _ := os.ReadFile(path)
			source <- content
		}
	}, func(item []byte, writer mapreduce.Writer[[]*types.KeyValue], cancel func(error)) {
		res := make([]*types.KeyValue, 0, 1e3)
		lines := strings.Split(string(item), "\r\n")
		for _, line := range lines[1:] {
			docStruct, _ := input_data.Doc2Struct(line)
			if docStruct.DocId == 0 {
				continue
			}
			// 分词
			tokens, _ := analyzer.GseCutForBuildIndex(docStruct.DocId, docStruct.Body)
			for _, v := range tokens {
				res = append(res, &types.KeyValue{Key: v.Token, Value: cast.ToString(v.DocId)})
			}
		}
		// shuffle 排序过程
		sort.Sort(types.ByKey(res))
		writer.Write(res)
	}, func(pipe <-chan []*types.KeyValue, writer mapreduce.Writer[string], cancel func(error)) {
		for values := range pipe {
			for _, v := range values { // 构建倒排索引
				if value, ok := invertedIndex.Get(v.Key); ok {
					value.AddInt(cast.ToInt(v.Value))
					invertedIndex.Set(v.Key, value)
				} else {
					docIds := roaring.NewBitmap()
					docIds.AddInt(cast.ToInt(v.Value))
					invertedIndex.Set(v.Key, docIds)
				}
			}
		}
	})

	keys := invertedIndex.Keys()
	for _, v := range keys {
		val, _ := invertedIndex.Get(v)
		fmt.Println(v, val)
	}

	return
}
