package job

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/RoaringBitmap/roaring"
	"github.com/golang-module/carbon"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/spf13/cast"

	"github.com/CocaineCong/tangseng/app/index_platform/consts"
	"github.com/CocaineCong/tangseng/app/index_platform/repository/storage"
	cconsts "github.com/CocaineCong/tangseng/consts"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/pkg/mapreduce"
	"github.com/CocaineCong/tangseng/pkg/timeutils"
	"github.com/CocaineCong/tangseng/repository/redis"
	"github.com/CocaineCong/tangseng/types"
)

// TODO:还有trie tree的merge,有空再补上

// MergeInvertedIndexDay2Month 增量合并全量, 合并到这个月，这周的数据就会删掉，下一个周增量又开始了
func MergeInvertedIndexDay2Month(ctx context.Context) (err error) {
	invertedIndexDayKey := redis.InvertedIndexDbPathDayKey
	invertedIndexMonthKey := redis.InvertedIndexDbPathMonthKey

	fromPaths, err := redis.ListInvertedPath(ctx, []string{invertedIndexDayKey})
	if err != nil {
		logs.LogrusObj.Errorln(err)
	}

	return mergeInvertedIndex(ctx, []string{invertedIndexDayKey}, fromPaths, invertedIndexMonthKey, consts.MergeTypeInvertedIndexDay2Month)
}

// MergeInvertedIndexMonth2Season 增量合并全量, 合并完就会删掉原有的，合并到这个季度
func MergeInvertedIndexMonth2Season(ctx context.Context) (err error) {
	invertedIndexMonthKey := redis.GetInvertedIndexDbPathMonthKey("*")
	invertedIndexSeasonKey := redis.InvertedIndexDbPathSeasonKey
	monthKeys, err := redis.ListAllPrefixKey(ctx, invertedIndexMonthKey)
	if err != nil {
		logs.LogrusObj.Error(err)
	}
	// 获取所有的月份的key
	fromPaths, err := redis.ListInvertedIndexByPrefixKey(ctx, invertedIndexMonthKey)
	if err != nil {
		logs.LogrusObj.Errorln(err)
	}

	return mergeInvertedIndex(ctx, monthKeys, fromPaths, invertedIndexSeasonKey, consts.MergeTypeInvertedIndexMonth2Season)
}

// mergeInvertedIndex fromPathKeys 所需要合并的key, fromPaths 需要合并的所有地址(就是key对应的地址)，toPathKey 合并完之后的存储该地址的key，mergeType，合并类型
func mergeInvertedIndex(ctx context.Context, fromPathKeys, fromPaths []string, savePathKey string, mergeType int) (err error) {
	invertedIndex := cmap.New[*roaring.Bitmap]() // 倒排索引
	_, _ = mapreduce.MapReduce(func(source chan<- []*types.InvertedInfo) {
		// 获取所有的inverted db
		for _, path := range fromPaths {
			invertedDb := storage.NewInvertedDB(path)
			p, _ := invertedDb.GetAllInverted()
			source <- p
		}
	}, func(item []*types.InvertedInfo, writer mapreduce.Writer[[]*types.InvertedInfo], cancel func(err error)) {
		// 对所有的inverted index进行对比
		sort.Slice(item, func(i, j int) bool {
			return item[i].Token < item[j].Token
		})
		writer.Write(item)
	}, func(pipe <-chan []*types.InvertedInfo, writer mapreduce.Writer[[]*types.InvertedInfo], cancel func(err error)) {
		// 整合所有的inverted index
		for values := range pipe {
			for _, v := range values { // 构建倒排索引
				if value, ok := invertedIndex.Get(v.Token); ok {
					value.AndAny(v.DocIds)
				} else {
					docIds := roaring.NewBitmap()
					docIds.AndAny(v.DocIds)
				}
			}
		}
	})

	// 生成所需要存储的key
	storageBaseName := ""
	switch mergeType {
	case consts.MergeTypeInvertedIndexDay2Month:
		storageBaseName = timeutils.GetMonthDate()
		savePathKey = redis.GetInvertedIndexDbPathMonthKey(cast.ToString(carbon.Now().Month()))
	case consts.MergeTypeInvertedIndexMonth2Season:
		storageBaseName = timeutils.GetSeasonDate()
		savePathKey = redis.GetInvertedIndexDbPathSeasonKey(cast.ToString(carbon.Now().Season()))
	default:
		storageBaseName = consts.InvertedIndexDefaultName
	}

	dir, _ := os.Getwd()
	outName := fmt.Sprintf("%s/%s.%s", dir, storageBaseName, cconsts.InvertedBucket)
	invertedDB := storage.NewInvertedDB(outName)
	// 找出所有的key进行存储
	for k, val := range invertedIndex.Items() {
		outByte, errx := val.MarshalBinary()
		if errx != nil {
			logs.LogrusObj.Error("mergeInvertedIndex-MarshalBinary", errx)
			continue
		}
		err = invertedDB.StoragePostings(k, outByte)
		if err != nil {
			logs.LogrusObj.Error("mergeInvertedIndex-StoragePostings", err)
			continue
		}
	}
	invertedDB.Close()

	// 保存新生成的索引数据地址
	err = redis.SetInvertedPath(ctx, savePathKey, outName)
	if err != nil {
		logs.LogrusObj.Error(err)
		return
	}

	// 删除旧纬度数据
	err = redis.BatchDeleteInvertedIndexPath(ctx, fromPathKeys)
	if err != nil {
		logs.LogrusObj.Error(err)
		return
	}

	return
}
