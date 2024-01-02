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

package recall

import (
	"context"
	"sort"

	"github.com/pkg/errors"

	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"
	"github.com/spf13/cast"

	"github.com/CocaineCong/tangseng/app/search_engine/analyzer"
	"github.com/CocaineCong/tangseng/app/search_engine/ranking"
	"github.com/CocaineCong/tangseng/app/search_engine/repository/db/dao"
	"github.com/CocaineCong/tangseng/app/search_engine/repository/storage"
	"github.com/CocaineCong/tangseng/app/search_engine/rpc"
	pb "github.com/CocaineCong/tangseng/idl/pb/search_vector"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/repository/redis"
	"github.com/CocaineCong/tangseng/types"
)

// Recall 查询召回
type Recall struct {
}

func NewRecall() *Recall {
	return &Recall{}
}

// Search 入口
func (r *Recall) Search(ctx context.Context, query string) (resp []*types.SearchItem, err error) {
	splitQuery, err := analyzer.GseCutForRecall(query)
	if err != nil {
		err = errors.WithMessagef(err, "text2postingslists error")
		return
	}

	// 倒排库搜索
	res, err := r.searchDoc(ctx, splitQuery)
	if err != nil {
		err = errors.WithMessage(err, "searchDoc error")
		return
	}

	// 向量库搜索
	vRes, err := r.SearchVector(ctx, splitQuery)
	if err != nil {
		err = errors.WithMessage(err, "searchVector error")
		return
	}

	resp, _ = r.Multiplex(ctx, query, res, vRes)
	return
}

// Multiplex 多路融合排序
func (r *Recall) Multiplex(ctx context.Context, query string, iRes, vRes []int64) (resp []*types.SearchItem, err error) {
	// 融合去重
	iRes = append(iRes, vRes...)
	iRes = lo.Uniq(iRes)
	recallData, _ := dao.NewInputDataDao(ctx).ListInputDataByDocIds(iRes)
	searchItems := make([]*types.SearchItem, 0)

	// 处理
	for _, v := range recallData {
		if v.Content == "" || v.Title == "" {
			continue
		}
		searchItems = append(searchItems, v)
	}

	// 排序
	searchItems, _ = ranking.CalculateScoreTFIDF(query, searchItems)
	sort.Slice(searchItems, func(i, j int) bool {
		return searchItems[i].Score > searchItems[j].Score
	})

	// 二次处理
	for _, v := range searchItems {
		if v.Score == 0 {
			continue
		}
		resp = append(resp, v)
	}

	return
}

// SearchVector 搜索向量
func (r *Recall) SearchVector(ctx context.Context, queries []string) (docIds []int64, err error) {
	// rpc 调用python接口 获取
	req := &pb.SearchVectorRequest{Query: queries}
	vectorResp, err := rpc.SearchVector(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "searchVector error")
		return
	}
	docIds = make([]int64, len(vectorResp.DocIds))
	for i, v := range vectorResp.DocIds {
		docIds[i] = cast.ToInt64(v)
	}
	return
}

// SearchQueryWord 入口词语联想
func (r *Recall) SearchQueryWord(query string) (resp []string, err error) {
	dictTreeList := make([]string, 0, 1e3)
	for _, trieDb := range storage.GlobalTrieDB {
		trie, errx := trieDb.GetTrieTreeDict()
		if errx != nil {
			err = errors.WithMessage(errx, "GetTrieTreeDict error")
			continue
		}
		queryTrie := trie.FindAllByPrefixForRecall(query)
		dictTreeList = append(dictTreeList, queryTrie...)
	}

	resp = lo.Uniq(dictTreeList)
	return
}

func (r *Recall) searchDoc(ctx context.Context, tokens []string) (recalls []int64, err error) {
	recalls = make([]int64, 0)
	for _, token := range tokens {
		docIds, errx := redis.GetInvertedIndexTokenDocIds(ctx, token)
		var postingsList []*types.PostingsList
		if errx != nil || docIds == nil {
			docIds = roaring.New()
			// 如果缓存不存在，就去索引表里面读取
			postingsList, err = fetchPostingsByToken(token)
			if err != nil {
				err = errors.WithMessage(err, "fetchPostingsByToken error")
				continue
			} else {
				// 如果缓存存在，就直接读缓存，不用担心实时性问题，缓存10分钟清空一次，这延迟是能接受到
				for _, v := range postingsList {
					if v != nil && v.DocIds != nil {
						docIds.AddMany(v.DocIds.ToArray())
					}
				}
			}

		}
		// TODO: term的position，后面再更新
		if docIds != nil {
			for _, v := range docIds.ToArray() {
				recalls = append(recalls, cast.ToInt64(v))
			}
		}
	}

	// 排序打分
	// iDao := dao.NewInputDataDao(ctx)
	// for _, p := range allPostingsList {
	// 	if p == nil || p.DocIds == nil || p.DocIds.IsEmpty() {
	// 		continue
	// 	}
	// 	recallData, _ := iDao.ListInputDataByDocIds(p.DocIds.ToArray())
	// 	searchItems := ranking.CalculateScoreBm25(p.Term, recallData)
	// 	recalls = append(recalls, searchItems...)
	// }

	log.LogrusObj.Infof("recalls size:%v", len(recalls))

	return
}

// 获取 token 所有seg的倒排表数据
func fetchPostingsByToken(token string) (postingsList []*types.PostingsList, err error) {
	// 遍历存储index的地方，token对应的doc Id 全部取出
	postingsList = make([]*types.PostingsList, 0, 1e6)
	for _, inverted := range storage.GlobalInvertedDB {
		docIds, errx := inverted.GetInverted([]byte(token))
		if errx != nil {
			err = errors.WithMessage(err, "getInverted error")
			continue
		}
		output := roaring.New()
		_ = output.UnmarshalBinary(docIds)
		// 存放到数组当中
		postings := &types.PostingsList{
			Term:   token,
			DocIds: output,
		}
		postingsList = append(postingsList, postings)
	}

	return
}
