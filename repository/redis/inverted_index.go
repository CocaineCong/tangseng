// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package redis

import (
	"context"

	"github.com/pkg/errors"

	"github.com/RoaringBitmap/roaring"
	"github.com/redis/go-redis/v9"
)

// PushInvertedPath 把存放db的path信息放到redis中
func PushInvertedPath(ctx context.Context, key string, paths []string) (err error) {
	for _, v := range paths {
		err = RedisClient.LPush(ctx, key, v).Err()
		if err != nil {
			return errors.Wrap(err, "failed to push inverted path in redis")
		}
	}

	return
}

// SetInvertedPath 把存放db的path信息放到redis中
func SetInvertedPath(ctx context.Context, key string, path string) (err error) {
	err = RedisClient.Set(ctx, key, path, redis.KeepTTL).Err()
	if err != nil {
		return errors.Wrap(err, "failed to set inverted path in redis")
	}

	return
}

// GetInvertedPath 获取存储的path信息
func GetInvertedPath(ctx context.Context, key string) (path string, err error) {
	path, err = RedisClient.Get(ctx, key).Result()
	if err != nil {
		return path, errors.Wrap(err, "failed to get inverted path")
	}

	return
}

// ListInvertedPath 把存放在redis的信息放到path中 包括day，week，season的
func ListInvertedPath(ctx context.Context, keys []string) (paths []string, err error) {
	for _, key := range keys {
		switch key {
		case InvertedIndexDbPathDayKey, TireTreeDbPathDayKey:
			results := RedisClient.LRange(ctx, key, 0, -1)
			paths = append(paths, results.Val()...)
		case InvertedIndexDbPathMonthKey, InvertedIndexDbPathSeasonKey,
			TireTreeDbPathMonthKey, TireTreeDbPathSeasonKey:
			// 由于这些key都不是列表，所以获取通过前缀获取所有的信息
			prefixKey := getAllDbPaths(key)
			results, errx := ListInvertedIndexByPrefixKey(ctx, prefixKey)
			if errx != nil {
				return paths, errors.Wrap(errx, "failed to list inverted index")
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
	err = RedisClient.Del(ctx, key).Err()
	if err != nil {
		return errors.Wrap(err, "failed to delete inverted index path")
	}
	return
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
	err = RedisClient.Set(ctx, getQueryTokenDocIdsKey(token), docIdsByte, QueryTokenDocIdsDefaultTimeout).Err()
	if err != nil {
		return errors.Wrap(err, "failed to set inverted index token docIds")
	}
	return
}

// GetInvertedIndexTokenDocIds 获取缓存的结果
func GetInvertedIndexTokenDocIds(ctx context.Context, token string) (docIds *roaring.Bitmap, err error) {
	res, err := RedisClient.Get(ctx, getQueryTokenDocIdsKey(token)).Result()
	if err != nil {
		return docIds, errors.Wrap(err, "failed to get query token docIds key from Redis")
	}
	docIds = roaring.NewBitmap()
	err = docIds.UnmarshalBinary([]byte(res))
	if err != nil {
		return docIds, errors.Wrap(err, "failed to unmarshal binary")
	}

	return
}

// PushInvertedIndexToken 存储用户搜索的历史记录 docs ids // TODO:后面嵌入LRU
func PushInvertedIndexToken(ctx context.Context, userId int64, token string) (err error) {
	err = RedisClient.LPush(ctx, getUserQueryTokenKey(userId), token).Err()
	if err != nil {
		return errors.Wrap(err, "failed to push inverted index token")
	}
	return
}

// ListInvertedIndexToken 获取用户搜索的历史记录
func ListInvertedIndexToken(ctx context.Context, userId int64) (tokens []string, err error) {
	tokens, err = RedisClient.LRange(ctx, getUserQueryTokenKey(userId), 0, -1).Result()
	if err != nil {
		return tokens, errors.Wrap(err, "failed to list inverted index token")
	}

	return
}

// PushInvertedMonthPath 把存放db的path信息放到redis中 month 纬度
func PushInvertedMonthPath(ctx context.Context, key string, paths []string) (err error) {
	for _, v := range paths {
		err = RedisClient.LPush(ctx, key, v).Err()
		if err != nil {
			return errors.Wrap(err, "failed to push inverted month path")
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
