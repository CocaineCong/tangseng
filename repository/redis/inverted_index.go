package redis

import (
	"context"

	"github.com/RoaringBitmap/roaring"
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

// ListInvertedPath 把存放在redis的信息放到path中
func ListInvertedPath(ctx context.Context, key string) (paths []string, err error) {
	results := RedisClient.LRange(ctx, key, 0, -1)
	if err != nil {
		return
	}
	paths = results.Val()

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
