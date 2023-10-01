package service

import (
	"context"
	"fmt"
	"hash/fnv"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/RoaringBitmap/roaring"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/spf13/cast"

	"github.com/CocaineCong/tangseng/app/index_platform/analyzer"
	"github.com/CocaineCong/tangseng/app/index_platform/input_data"
	"github.com/CocaineCong/tangseng/app/index_platform/repository/storage"
	cconsts "github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/consts/e"
	pb "github.com/CocaineCong/tangseng/idl/pb/index_platform"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/pkg/mapreduce"
	"github.com/CocaineCong/tangseng/pkg/timeutils"
	"github.com/CocaineCong/tangseng/pkg/trie"
	"github.com/CocaineCong/tangseng/repository/redis"
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
	dictTrie := trie.NewTrie()                   // 前缀树

	logs.LogrusObj.Infof("BuildIndexService Start req: %v", req.FilePath)
	// mapreduce 这个是用chan和goroutine来代替master和worker的rpc调用，避免了频繁的rpc调用
	_, _ = mapreduce.MapReduce(func(source chan<- []byte) {
		for _, path := range req.FilePath {
			content, _ := os.ReadFile(path)
			source <- content
		}
	}, func(item []byte, writer mapreduce.Writer[[]*types.KeyValue], cancel func(error)) {
		// 控制并发
		var wg sync.WaitGroup
		ch := make(chan struct{}, 3)

		keyValueList := make([]*types.KeyValue, 0, 1e3)
		lines := strings.Split(string(item), "\r\n")
		for _, line := range lines[1:] {
			ch <- struct{}{}
			wg.Add(1)
			docStruct, _ := input_data.Doc2Struct(line) // line 转 doc struct
			if docStruct.DocId == 0 {
				continue
			}

			// 分词
			tokens, _ := analyzer.GseCutForBuildIndex(docStruct.DocId, docStruct.Body)
			for _, v := range tokens {
				if v.Token == "" || v.Token == " " {
					continue
				}
				keyValueList = append(keyValueList, &types.KeyValue{Key: v.Token, Value: cast.ToString(v.DocId)})
				dictTrie.Insert(v.Token)
			}

			// 建立正排索引
			go func(docStruct *types.Document) {
				err = input_data.DocData2Kfk(docStruct)
				if err != nil {
					logs.LogrusObj.Error(err)
				}
				defer wg.Done()
				<-ch
			}(docStruct)
		}
		wg.Wait()

		// // 构建前缀树 // TODO:kafka处理
		// go func(tokenList []string) {
		// 	err = input_data.DocTrie2Kfk(tokenList)
		// 	if err != nil {
		// 		logs.LogrusObj.Error("DocTrie2Kfk", err)
		// 	}
		// }(tokenList)

		// shuffle 排序过程
		sort.Sort(types.ByKey(keyValueList))
		writer.Write(keyValueList)
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

	go func() {
		err = storeInvertedIndexByHash(ctx, invertedIndex)
		if err != nil {
			logs.LogrusObj.Error("storeInvertedIndexByHash error ", err)
		}
	}()

	logs.LogrusObj.Infoln("storeInvertedIndexByHash End")
	go func() {
		err = storeDictTrieByHash(ctx, dictTrie)
		if err != nil {
			logs.LogrusObj.Error("storeDictTrieByHash error ", err)
		}
	}()

	return
}

// storeInvertedIndexByHash 分片存储
func storeInvertedIndexByHash(ctx context.Context, invertedIndex cmap.ConcurrentMap[string, *roaring.Bitmap]) (err error) {
	dir, _ := os.Getwd()
	outName := fmt.Sprintf("%s/%s.%s", dir, timeutils.GetTodayDate(), cconsts.InvertedBucket)
	invertedDB := storage.NewInvertedDB(outName)
	// 找出所有的key进行存储
	for k, val := range invertedIndex.Items() {
		outByte, errx := val.MarshalBinary()
		if errx != nil {
			logs.LogrusObj.Error("storeInvertedIndexByHash-MarshalBinary", errx)
			continue
		}
		err = invertedDB.StoragePostings(k, outByte)
		if err != nil {
			logs.LogrusObj.Error("storeInvertedIndexByHash-StoragePostings", err)
			continue
		}
	}

	err = redis.PushInvertedPath(ctx, redis.InvertedIndexDbPathKey, []string{outName})
	if err != nil {
		logs.LogrusObj.Error(err)
		return
	}

	// TODO: hash 分片存储
	// dir, _ := os.Getwd()
	// keys := invertedIndex.Keys()
	// buffer := make([][]*types.KeyValue, consts.ReduceDefaultNum)
	// for i, v := range keys {
	// 	val, _ := invertedIndex.Get(v)
	// 	slot := iHash(v) % consts.ReduceDefaultNum
	// 	buffer[slot] = append(buffer[slot])
	// 	fmt.Println(v, val)
	// }
	// outName := fmt.Sprintf("%s/%d.%s", dir, i, cconsts.InvertedBucket)

	return
}

// storeInvertedIndexByHash 分片存储
func storeDictTrieByHash(ctx context.Context, dict *trie.Trie) (err error) {
	// TODO: 抽离一个hash存储的方法
	dir, _ := os.Getwd()
	outName := fmt.Sprintf("%s/%s.%s", dir, timeutils.GetTodayDate(), cconsts.TrieTreeBucket)
	trieDB := storage.NewTrieDB(outName)
	err = trieDB.StorageDict(dict)
	if err != nil {
		logs.LogrusObj.Error(err)
		return
	}

	err = redis.PushInvertedPath(ctx, redis.TireTreeDbPathKey, []string{outName})
	if err != nil {
		logs.LogrusObj.Error(err)
		return
	}

	return
}

func iHash(key string) int64 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return int64(h.Sum32() & 0x7fffffff)
}
