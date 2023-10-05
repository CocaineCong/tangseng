package redis

import (
	"context"

	"github.com/RoaringBitmap/roaring"
	"github.com/redis/go-redis/v9"
)

// PushInvertedPath 把存放db的path信息放到redis中
func PushInvertedPath(ctx context.Context, key string, paths []string) (err error) {
	for _, v := range paths {
		err = RedisClient.LPush(ctx, key, v).Err()
		if err != nil {
			return err
		}
	}

	return
}

// SetInvertedPath 把存放db的path信息放到redis中
func SetInvertedPath(ctx context.Context, key string, path string) (err error) {
	err = RedisClient.Set(ctx, key, path, redis.KeepTTL).Err()
	if err != nil {
		return err
	}

	return
}

// GetInvertedPath 获取存储的path信息
func GetInvertedPath(ctx context.Context, key string) (path string, err error) {
	path, err = RedisClient.Get(ctx, key).Result()
	if err != nil {
		return
	}

	return
}

// ListInvertedPath 把存放在redis的信息放到path中 包括day，week，season的
func ListInvertedPath(ctx context.Context, keys []string) (paths []string, err error) {
	for _, key := range keys {
		switch key {
		case InvertedIndexDbPathDayKey, TireTreeDbPathDayKey:
			results := RedisClient.LRange(ctx, key, 0, -1)
			if err != nil {
				return
			}
			paths = append(paths, results.Val()...)
		case InvertedIndexDbPathMonthKey, InvertedIndexDbPathSeasonKey,
			TireTreeDbPathMonthKey, TireTreeDbPathSeasonKey:
			// 由于这些key都不是列表，所以获取通过前缀获取所有的信息
			prefixKey := getAllDbPaths(key)
			results, errx := ListInvertedIndexByPrefixKey(ctx, prefixKey)
			if errx != nil {
				err = errx
				return
			}
			paths = append(paths, results...)
		default:
			return
		}
	}

	return
}

// DeleteInvertedIndexPath 删除 inverted index path
func DeleteInvertedIndexPath(ctx context.Context, key string) (err error) {
	return RedisClient.Del(ctx, key).Err()
}

// BatchDeleteInvertedIndexPath 批量删除 inverted index path
func BatchDeleteInvertedIndexPath(ctx context.Context, keys []string) (err error) {
	for _, key := range keys {
		_ = DeleteInvertedIndexPath(ctx, key)
	}

	return
}

// SetInvertedIndexTokenDocIds 缓存搜索过的结果 // TODO:后面嵌入LRU
func SetInvertedIndexTokenDocIds(ctx context.Context, token string, docIds *roaring.Bitmap) (err error) {
	docIdsByte, _ := docIds.MarshalBinary()
	return RedisClient.Set(ctx, getQueryTokenDocIdsKey(token), docIdsByte, QueryTokenDocIdsDefaultTimeout).Err()
}

// GetInvertedIndexTokenDocIds 获取缓存的结果
func GetInvertedIndexTokenDocIds(ctx context.Context, token string) (docIds *roaring.Bitmap, err error) {
	res, err := RedisClient.Get(ctx, getQueryTokenDocIdsKey(token)).Result()
	if err != nil {
		return
	}
	docIds = roaring.NewBitmap()
	err = docIds.UnmarshalBinary([]byte(res))
	if err != nil {
		return
	}

	return
}

// PushInvertedIndexToken 存储用户搜索的历史记录 docs ids // TODO:后面嵌入LRU
func PushInvertedIndexToken(ctx context.Context, userId int64, token string) (err error) {
	return RedisClient.LPush(ctx, getUserQueryTokenKey(userId), token).Err()
}

// ListInvertedIndexToken 获取用户搜索的历史记录
func ListInvertedIndexToken(ctx context.Context, userId int64) (tokens []string, err error) {
	tokens, err = RedisClient.LRange(ctx, getUserQueryTokenKey(userId), 0, -1).Result()
	if err != nil {
		return
	}

	return
}

// PushInvertedMonthPath 把存放db的path信息放到redis中 month 纬度
func PushInvertedMonthPath(ctx context.Context, key string, paths []string) (err error) {
	for _, v := range paths {
		err = RedisClient.LPush(ctx, key, v).Err()
		if err != nil {
			return err
		}
	}

	return
}

// ListInvertedIndexByPrefixKey 通过前缀获取所有的value index_platform:inverted_index:month:*
func ListInvertedIndexByPrefixKey(ctx context.Context, prefixKey string) (paths []string, err error) {
	// 使用Scan方法遍历所有键
	iter := RedisClient.Scan(ctx, 0, prefixKey, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		value, _ := RedisClient.Get(ctx, key).Result()
		paths = append(paths, value)
	}

	return
}

// ListAllPrefixKey 通过前缀获取所有的value index_platform:inverted_index:month:*
func ListAllPrefixKey(ctx context.Context, prefixKey string) (paths []string, err error) {
	// 使用Scan方法遍历所有键
	iter := RedisClient.Scan(ctx, 0, prefixKey, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		paths = append(paths, key)
	}

	return
}
